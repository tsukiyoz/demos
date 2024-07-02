package v1

import (
	"fmt"
	"log"
)

type Worker struct {
	id   int
	jobs *JobQueue
}

func NewWorker(id int, jobs *JobQueue) *Worker {
	return &Worker{
		id:   id,
		jobs: jobs,
	}
}

func (w *Worker) handleCrash() {
	r := recover()
	if r != nil {
		log.Printf("worker %d, recovered form panic: %#v", w.id, r)
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			select {
			case <-w.jobs.wait():
				fmt.Printf("worker: %d, get a job\n", w.id)
				job := w.jobs.Pop()
				fmt.Printf("worker: %d, execute job started\n", w.id)
				func() {
					defer w.handleCrash()
					_ = job.Execute()
				}()
				fmt.Printf("worker: %d, execute job done\n", w.id)
				job.Done()
			}
		}
	}()
}
