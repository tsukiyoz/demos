package taskschedule

import (
	"errors"
	"fmt"
)

type Scheduler[T Task] struct {
	taskCh       chan T
	execute      ExecuteStrategy[T]
	paralism     chan struct{}
	stateHandler map[TaskStatus]TaskHandler[T]
}

type Option[T Task] func(*Scheduler[T])

func WithTaskDefaultSize[T Task](size int) Option[T] {
	return func(s *Scheduler[T]) {
		s.taskCh = make(chan T, size)
	}
}

func WithExecuteStrategy[T Task](es ExecuteStrategy[T]) Option[T] {
	return func(s *Scheduler[T]) {
		s.execute = es
	}
}

func WithParallelism[T Task](size int) Option[T] {
	return func(s *Scheduler[T]) {
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

type TaskHandler[T Task] func(T) error

func NewScheduler[T Task](handler map[TaskStatus]TaskHandler[T], opts ...Option[T]) *Scheduler[T] {
	s := &Scheduler[T]{
		taskCh:       make(chan T, defaultTaskSize),
		execute:      NewParallelStrategy[T](),
		paralism:     make(chan struct{}, defaultParallel),
		stateHandler: handler,
	}
	for _, opt := range opts {
		opt(s)
	}
	go s.run()
	return s
}

func (s *Scheduler[T]) AddTask(task T) {
	s.taskCh <- task
}

func (s *Scheduler[T]) Close() {
	close(s.taskCh)
	ps := len(s.paralism)
	for i := 0; i < ps; i++ {
		s.paralism <- struct{}{}
	}
}

func (s *Scheduler[T]) run() {
	for task := range s.taskCh {
		s.paralism <- struct{}{}

		go func() {
			defer func() {
				<-s.paralism
			}()
			if err := s.execute(s, task); err != nil {
				fmt.Printf("task %d failed: %v\n", task.ID(), err)
				return
			}
			if callback := task.Callback(); callback != nil {
				if err := callback(task); err != nil {
					fmt.Printf("task %d callback failed: %v\n", task.ID(), err)
				}
			}
		}()
	}
}

type ExecuteStrategy[T Task] func(s *Scheduler[T], t T) error

func NewFallthroughStrategy[T Task]() ExecuteStrategy[T] {
	return func(s *Scheduler[T], t T) error {
		for t.Next() {
			handler, ok := s.stateHandler[t.Status()]
			if !ok {
				return errors.New("unknown status")
			}

			if err := handler(t); err != nil {
				return err
			}
		}
		return nil
	}
}

func NewParallelStrategy[T Task]() ExecuteStrategy[T] {
	return func(s *Scheduler[T], t T) error {
		handler, ok := s.stateHandler[t.Status()]
		if !ok {
			return errors.New("unknown status")
		}
		err := handler(t)
		if err != nil {
			return err
		}
		if t.Next() {
			go func() { s.AddTask(t) }()
		}
		return nil
	}
}
