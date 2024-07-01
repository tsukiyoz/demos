package http

import (
	"container/list"
	"fmt"
	"sync"
)

type JobQueue struct {
	mu       sync.Mutex
	notice   chan signal
	queue    *list.List
	size     int
	capacity int
}

func NewJobQueue(cap int) *JobQueue {
	return &JobQueue{
		capacity: cap,
		queue:    list.New(),
		notice:   make(chan signal, 1),
	}
}

func (q *JobQueue) Push(j *Job) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.size >= q.capacity {
		q.RemoveLeastJob()
	}
	q.size++
	q.queue.PushBack(j)
	q.notice <- signal{}
}

func (q *JobQueue) RemoveLeastJob() {
	if q.queue.Len() != 0 {
		back := q.queue.Back()
		job := back.Value.(*Job)
		job.Done()
		fmt.Printf("kill a job\n")
		q.queue.Remove(back)
		q.size--
	}
}

func (q *JobQueue) Pop() *Job {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.size == 0 {
		return nil
	}

	q.size--
	return q.queue.Remove(q.queue.Front()).(*Job)
}

func (q *JobQueue) wait() <-chan signal {
	return q.notice
}
