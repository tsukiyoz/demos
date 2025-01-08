package timewheel_test

import (
	"sync"
	"testing"
	"time"

	"github.com/lazywoo/demos/algo/timewheel"
)

func TestAdd(t *testing.T) {
	tw := timewheel.New(time.Second, 3600)
	tw.Start()

	wg := sync.WaitGroup{}

	wg.Add(2)

	tw.Add(time.Second*5, "", func() {
		t.Log("empty key task executed in 5 seconds")
		wg.Done()
	})

	tw.Add(time.Second*7, "task", func() {
		t.Log("task executed in 7 seconds")
		wg.Done()
	})

	wg.Wait()
}

func TestCancel(t *testing.T) {
	tw := timewheel.New(time.Second, 3600)
	tw.Start()

	wg := sync.WaitGroup{}

	wg.Add(2)

	tw.Add(time.Second*5, "", func() {
		t.Log("empty key task executed in 5 seconds")
		wg.Done()
	})

	tw.Add(time.Second*7, "task", func() {
		t.Log("task executed in 7 seconds")
		wg.Done()
	})

	tw.Cancel("task")
	wg.Done()

	wg.Wait()
}
