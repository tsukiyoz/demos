package v1

type signal struct{}

type Job struct {
	done       chan signal
	handleFunc func() error
}

func (j *Job) Execute() error {
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
