package main

import (
	"fmt"
	"time"
)

type Large [1 << 12]byte

func foo() {
	for a, i := (Large{}), 0; i < len(a); i++ {
		f(&a, i)
	}
}

func f(x *Large, k int) {}

func main() {
	bench := func() time.Duration {
		start := time.Now()
		foo()
		return time.Since(start)
	}
	fmt.Println("duration: ", bench())
}
