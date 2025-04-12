package main

// go build -gcflags="-m -l" ./escape.go

func stackAllocated(n int) {
	a := make([]int, 0, 10) // make([]int, 0, 10) does not escape
	a = append(a, 1)
	_ = a
}

func heapAllcated(n int) []int {
	a := make([]int, 0, n) // make([]int, 0, n) escapes to heap
	a = append(a, 1)
	return a
}

func closureAllocated(n int) func() []int {
	a := make([]int, 0, n) // make([]int, 0, n) escapes to heap
	return func() []int {
		a = append(a, 1)
		return a
	}
}

func main() {
	n := 20
	stackAllocated(n)
	_ = heapAllcated(n)
	_ = closureAllocated(n)
}
