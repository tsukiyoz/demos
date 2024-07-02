package v2

import (
	"fmt"
	"log"
)

type Worker struct {
	id int
	q  chan *Job
}

func NewWorker(id int, q chan *Job) *Worker {
	wk := &Worker{
		id: id,
		q:  q,
	}
	go wk.run()
	return wk
}

func (wk *Worker) handleCrash() {
	r := recover()
	if r != nil {
		log.Printf("recovered form panic")
	}
}

func (wk *Worker) run() {
	fmt.Printf("worker %d is running\n", wk.id)
	for {
		select {
		case j := <-wk.q:
			fmt.Printf("worker %d get a job\n", wk.id)
			func() {
				defer wk.handleCrash()
				_ = wk.exec(j)
			}()
			j.Done()
			fmt.Printf("worker %d finished a job\n", wk.id)
		}
	}
}

func (wk *Worker) exec(j *Job) error {
	return j.fn()
}
