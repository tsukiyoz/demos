package main

import (
	"context"
	"fmt"
	"github.com/dustin/go-broadcast"
	"sync"
	"testing"
	"time"
)

func TestLiveLock(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cadence := sync.NewCond(&sync.Mutex{})
	go func() {
		for range time.Tick(time.Second) {
			cadence.Broadcast()
		}
	}()

	var wg sync.WaitGroup
	workerNum := 10
	wg.Add(workerNum)

	t.Logf("start hard working, :( \n")

	for i := range workerNum {
		go func(id int) {
			defer wg.Done()
			for {
				cadence.L.Lock()
				cadence.Wait()
				cadence.L.Unlock()
				select {
				case <-ctx.Done():
					return
				default:
					t.Logf("worker id:%v working...\n", id)
				}
			}
		}(i)
	}

	wg.Wait()
	t.Logf("go off work, :) \n")
}

type empty struct{}

var workerNum = 4

func BenchmarkBroadcastByChan(b *testing.B) {

}

func TestBroadcastByChan(t *testing.T) {
	doBroadcastByChan()
}

func doBroadcastByChan() {
	broadcaster := broadcast.NewBroadcaster(100)

	var wg sync.WaitGroup
	wg.Add(workerNum)
	for i := 0; i < workerNum; i++ {
		startWorker(broadcaster, &wg)
	}

	broadcaster.Submit(empty{})
	wg.Wait()
	broadcaster.Close()
}

func startWorker(broadcaster broadcast.Broadcaster, wg *sync.WaitGroup) {
	ch := make(chan interface{}, 1)
	broadcaster.Register(ch)
	defer broadcaster.Unregister(ch)

	go func() {
		select {
		case <-ch:
			fmt.Println("got one")
			wg.Done()
			return
		}
	}()
}
