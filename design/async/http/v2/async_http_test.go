package v2

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

const (
	MaxProcessorNum = 16
)

func TestHTTPServer(t *testing.T) {
	buflen := 50000
	q := make(chan *Job, MaxProcessorNum*buflen)

	var workers []*Worker
	for i := range MaxProcessorNum * 10000 {
		workers = append(workers, NewWorker(i+1, q))
	}

	http.Handle("/", NewHandler(q))
	http.ListenAndServe(":8080", nil)
}

type Handler struct {
	q chan *Job
}

func NewHandler(q chan *Job) *Handler {
	return &Handler{
		q: q,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handle http request")
	j := &Job{
		done: make(chan signal),
		fn: func() error {
			w.WriteHeader(http.StatusOK)
			// mock execute complex operations
			// time.Sleep(time.Millisecond * time.Duration(rand.Intn(20)+10))
			time.Sleep(time.Millisecond * time.Duration(30))
			// io.WriteString(w, "hello world, async http")
			w.Write([]byte("hello world, async http"))
			return nil
		},
	}
	startAt := time.Now()
	// send
	go func() {
		select {
		case h.q <- j:
			return
		case <-time.After(time.Second):
			w.WriteHeader(http.StatusRequestTimeout)
		}
	}()
	j.Wait()
	fmt.Printf("duration: %v\n", time.Since(startAt))
}
