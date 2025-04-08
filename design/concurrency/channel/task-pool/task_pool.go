package taskpool

import (
	"context"
	"sync"
)

type Task func() error

type TaskPool struct {
	tasks      chan Task
	maxWorkers int

	closeOnce sync.Once
	quit      chan struct{}
	wg        sync.WaitGroup

	onError func(error)
}

func NewTaskPool(maxWorkers int, bufferSize int) *TaskPool {
	if maxWorkers <= 0 {
		maxWorkers = 1
	}
	if bufferSize <= 0 {
		bufferSize = maxWorkers * 2
	}

	taskPool := &TaskPool{
		tasks:      make(chan Task, bufferSize),
		maxWorkers: maxWorkers,
		quit:       make(chan struct{}),
	}

	taskPool.wg.Add(maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		go func() {
			defer taskPool.wg.Done()
			for {
				select {
				case task, ok := <-taskPool.tasks:
					if !ok {
						return
					}
					if err := task(); err != nil && taskPool.onError != nil {
						taskPool.onError(err)
					}
				case <-taskPool.quit:
					return
				}
			}
		}()
	}
	return taskPool
}

func (t *TaskPool) Submit(ctx context.Context, task Task) bool {
	select {
	case <-ctx.Done():
		// context is done, do not submit the task
		return false
	case <-t.quit:
		return false
	case t.tasks <- task:
		// task submitted successfully
		return true
	}
}

func (t *TaskPool) Close() error {
	t.closeOnce.Do(func() {
		close(t.quit)
		t.wg.Wait()
		close(t.tasks)
	})
	return nil
}

func (t *TaskPool) SetOnError(fn func(error)) {
	t.onError = fn
}
