package v1

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

type Handler struct {
	jobs *JobQueue
}

func NewHandler(jobs *JobQueue) *Handler {
	return &Handler{
		jobs: jobs,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	job := &Job{
		done: make(chan signal, 1),
		handleFunc: func() error {
			w.Write([]byte("hello world"))
			return nil
		},
	}
	h.jobs.Push(job)
	fmt.Println("commit job success")
	start := time.Now()
	job.Wait()
	fmt.Printf("duration: %v\n", time.Since(start))
}

func TestHTTPServer(t *testing.T) {
	jobs := NewJobQueue(24)

	handler := NewHandler(jobs)
	NewWorker(1, jobs).Start()
	NewWorker(2, jobs).Start()
	NewWorker(3, jobs).Start()
	NewWorker(4, jobs).Start()
	NewWorker(5, jobs).Start()

	http.Handle("/", handler)
	http.ListenAndServe(":8080", nil)
}
