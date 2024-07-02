package v2

import "sync"

type signal struct{}

type Job struct {
	sync.Once
	done chan signal
	fn   func() error
}

func (j *Job) Exec() error {
	return j.fn()
}

func (j *Job) Done() {
	j.Do(func() {
		j.done <- signal{}
	})
}

func (j *Job) Wait() {
	select {
	case <-j.done:
		return
	}
}
