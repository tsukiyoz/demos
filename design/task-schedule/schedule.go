package taskschedule

import (
	"errors"
)

type Scheduler[T Task] struct {
	taskCh       chan T
	execute      ExecuteStrategy[T]
	paralism     chan struct{}
	stateHandler map[TaskStatus]TaskHandler[T]
	afterExec    func(T)
	onError      func(T, error)
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

func WithAfterExec[T Task](fn func(T)) Option[T] {
	return func(s *Scheduler[T]) {
		s.afterExec = fn
	}
}

func WithOnError[T Task](fn func(T, error)) Option[T] {
	return func(s *Scheduler[T]) {
		s.onError = fn
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

func (s *Scheduler[T]) Submit(task T) {
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
			err := s.execute(s, task)
			if s.afterExec != nil {
				defer s.afterExec(task)
			}
			if err != nil && s.onError != nil {
				s.onError(task, err)
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
			go func() { s.Submit(t) }()
		}
		return nil
	}
}
