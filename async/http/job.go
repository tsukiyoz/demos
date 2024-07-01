package http

import "math/rand"

type signal struct{}

type Job struct {
	done       chan signal
	handleFunc func() error
}

func (j *Job) Execute() error {
	danger := rand.Intn(100)
	if danger > 80 {
		panic("danger operate!")
	}
	return j.handleFunc()
}

func (j *Job) Done() {
	j.done <- signal{}
	close(j.done)
}

func (j *Job) Wait() {
	select {
	case <-j.done:
		return
	}
}
