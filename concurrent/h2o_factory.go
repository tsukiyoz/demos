package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/marusama/cyclicbarrier"
	"golang.org/x/sync/semaphore"
)

type H2O struct {
	H *semaphore.Weighted
	O *semaphore.Weighted
	b cyclicbarrier.CyclicBarrier
}

func New() *H2O {
	return &H2O{
		H: semaphore.NewWeighted(2),
		O: semaphore.NewWeighted(1),
		b: cyclicbarrier.New(3),
	}
}

var ch chan string

func releaseHydrogen() {
	ch <- "H"
}

func releaseOxygen() {
	ch <- "O"
}

func (h2o *H2O) hydrogen(releaseHydrogen func()) {
	h2o.H.Acquire(context.Background(), 1)
	releaseHydrogen()
	h2o.b.Await(context.Background())
	h2o.H.Release(1)
}

func (h2o *H2O) oxygen(releaseOxygen func()) {
	h2o.O.Acquire(context.Background(), 1)
	releaseOxygen()
	h2o.b.Await(context.Background())
	h2o.O.Release(1)
}

func main() {
	N := 100
	ch = make(chan string, N*3)

	h2o := New()

	var wg sync.WaitGroup
	wg.Add(3 * N)

	for i := 0; i < 2*N; i++ {
		go func() {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			h2o.hydrogen(releaseHydrogen)
			wg.Done()
		}()
	}

	for i := 0; i < N; i++ {
		go func() {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			h2o.oxygen(releaseOxygen)
			wg.Done()
		}()
	}

	wg.Wait()

	if len(ch) != N*3 {
		fmt.Println(fmt.Errorf("expected %d items, got %d", N*3, len(ch)))
	}

	s := make([]string, 3)
	for i := 0; i < N; i++ {
		s[0] = <-ch
		s[1] = <-ch
		s[2] = <-ch
		// sort.Strings(s)

		water := s[0] + s[1] + s[2]
		//if water != "HHO" {
		//	fmt.Println(fmt.Errorf("expected HHO, got %s", water))
		//}

		fmt.Println(water)
	}
}
