package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
)

type (
	Value  string
	Signal struct{}
)

type Manager struct {
	wg      *sync.WaitGroup
	workers []*Worker
	inputs  chan empty
	results chan Value
	ctx     context.Context
	close   func()
	once    sync.Once
}

func (m *Manager) Start(ctx context.Context) {
	m.ctx, m.close = context.WithCancel(ctx)

	fmt.Printf("manager starting...\n")
	var wg sync.WaitGroup
	m.wg = &wg
	for _, worker := range m.workers {
		worker := worker
		go worker.start(&wg)
	}
	fmt.Printf("manager started\n")
	go m.healthz()
}

func (m *Manager) healthz() {
	for {
		select {
		case <-m.ctx.Done():
			_ = m.Close()
			return
		case <-time.After(time.Second * 3):
			fmt.Println("manager keep alive ...")
		}
	}
}

func (m *Manager) GetID() Value {
	m.inputs <- empty{}
	res := <-m.results
	return res
}

func (m *Manager) Close() error {
	m.once.Do(func() {
		fmt.Printf("manager exiting...\n")
		m.close()
		close(m.inputs)
		m.wg.Wait()
	})
	fmt.Printf("manager exited!\n")
	return nil
}

type Worker struct {
	ID   int
	mgr  *Manager
	data chan Value
}

func (w *Worker) start(wg *sync.WaitGroup) {
	defer wg.Done()

	wg.Add(1)
	fmt.Printf("worker[%d] started\n", w.ID)

	for range w.mgr.inputs {
		w.data <- Value(w.GetID())
	}

	fmt.Printf("worker[%d] quit....\n", w.ID)
}

func (w *Worker) GetID() string {
	fmt.Printf("worker[%d] working...\n", w.ID)
	return uuid.New().String()
}

func NewManager(workerNum int) *Manager {
	mgr := &Manager{
		inputs:  make(chan empty),
		results: make(chan Value),
	}
	workers := make([]*Worker, 0, workerNum)
	for i := range workerNum {
		workers = append(workers, &Worker{
			ID:   i + 1,
			mgr:  mgr,
			data: make(chan Value),
		})
	}
	mgr.workers = workers
	return mgr
}

func TestWorkersStart(t *testing.T) {
	mgr := NewManager(7)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	mgr.Start(ctx)

	reqN := 100
	for i := 0; i < reqN; i++ {
		t.Logf("get %v\n", mgr.GetID())
	}

	time.Sleep(time.Second * 16)
	mgr.Close()
}
