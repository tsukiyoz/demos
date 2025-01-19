package taskschedule

import (
	"errors"
	"fmt"
)

type Scheduler struct {
	taskCh   chan *Task
	execute  ExecuteStrategy
	paralism chan struct{}
}

type Option func(*Scheduler)

func WithTaskDefaultSize(size int) Option {
	return func(s *Scheduler) {
		s.taskCh = make(chan *Task, size)
	}
}

func WithExecuteStrategy(es ExecuteStrategy) Option {
	return func(s *Scheduler) {
		s.execute = es
	}
}

func WithParallelism(size int) Option {
	return func(s *Scheduler) {
		if size <= 0 {
			size = defaultParallel
		}
		s.paralism = make(chan struct{}, size)
	}
}

const (
	defaultTaskSize = 10
	defaultParallel = 10
)

func NewScheduler(opts ...Option) *Scheduler {
	s := &Scheduler{
		taskCh:   make(chan *Task, defaultTaskSize),
		execute:  NewParallelStrategy(),
		paralism: make(chan struct{}, defaultParallel),
	}
	for _, opt := range opts {
		opt(s)
	}
	go s.run()
	return s
}

func (s *Scheduler) AddTask(task *Task) {
	s.taskCh <- task
}

func (s *Scheduler) Close() {
	close(s.taskCh)
	ps := len(s.paralism)
	for i := 0; i < ps; i++ {
		s.paralism <- struct{}{}
	}
}

func (s *Scheduler) run() {
	for task := range s.taskCh {
		s.paralism <- struct{}{}

		go func() {
			defer func() {
				<-s.paralism
			}()
			if err := s.execute(s, task); err != nil {
				fmt.Printf("task %d failed: %v\n", task.ID, err)
				return
			}
			if task.Callback != nil {
				if err := task.Callback(task); err != nil {
					fmt.Printf("task %d callback failed: %v\n", task.ID, err)
				}
			}
		}()
	}
}

type ExecuteStrategy func(s *Scheduler, t *Task) error

func NewFallthroughStrategy() ExecuteStrategy {
	return func(s *Scheduler, t *Task) error {
		for {
			_, next := t.Status.Next()
			if !next {
				return nil
			}
			handler, ok := stateHandler[t.Status]
			if !ok {
				return errors.New("unknown status")
			}

			if err := handler(t); err != nil {
				return err
			}
		}
	}
}

func NewParallelStrategy() ExecuteStrategy {
	return func(s *Scheduler, t *Task) error {
		handler, ok := stateHandler[t.Status]
		if !ok {
			return errors.New("unknown status")
		}
		err := handler(t)
		if err != nil {
			return err
		}
		_, next := t.Status.Next()
		if !next {
			return nil
		}
		go func() { s.AddTask(t) }()
		return nil
	}
}
