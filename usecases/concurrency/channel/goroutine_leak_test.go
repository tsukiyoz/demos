package channel

import (
	"runtime"
	"sync"
	"testing"
)

func TestGoroutineLeak(t *testing.T) {
	memConsumed := func() uint64 {
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}

	var c <-chan struct{}
	var wg sync.WaitGroup
	noop := func() {
		wg.Done()
		<-c
	}
	before := memConsumed()
	const N = 1000
	wg.Add(N)
	for range N {
		go noop()
	}
	wg.Wait()

	after := memConsumed()
	t.Logf("cost %.3f kb", float64(after-before)/1024/float64(N)) // 5.184 kb
}
