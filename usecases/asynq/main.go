package main

import (
	"fmt"

	"golang.org/x/sync/errgroup"
)

func main() {
	var eg errgroup.Group

	eg.Go(startPublisher)
	eg.Go(startConsumer)

	err := eg.Wait()
	if err != nil {
		fmt.Println("error:", err)
	}
}
