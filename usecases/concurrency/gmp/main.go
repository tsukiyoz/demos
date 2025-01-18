package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	runtime.GOMAXPROCS(1)
	fmt.Println("Number of CPUs: ", runtime.NumCPU())
	fmt.Println("Number of Machines: ", runtime.GOMAXPROCS(0))

	printNumberAlternative()
}

func printNumberAlternative() {
	var wg sync.WaitGroup
	printFunc := func(s string, sched bool) {
		defer wg.Done()
		if sched {
			runtime.Gosched()
		}
		fmt.Println(s)
	}

	for i := 1; i <= 10; i++ {
		wg.Add(2)
		go printFunc("A", false)
		go printFunc("B", true)
		time.Sleep(100 * time.Millisecond)
	}

	wg.Wait()
}
