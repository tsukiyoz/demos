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
	q := make(chan *Job, MaxProcessorNum*4)

	var workers []*Worker
	for i := 0; i < MaxProcessorNum; i++ {
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

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handle http request")
	j := &Job{
		done: make(chan signal),
		fn: func() error {
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
