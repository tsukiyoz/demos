package test

import (
	"fmt"
	"log"
	"sync"
	"testing"
)

type WRRBalancer struct {
	nodes []*Node
	total int
	mu    sync.Mutex
}

func NewWRRBalancer(ws []int) *WRRBalancer {
	tot := 0
	nodes := make([]*Node, 0, len(ws))
	for i := 0; i < len(ws); i++ {
		tot += ws[i]
		nodes = append(nodes, &Node{
			name:   fmt.Sprintf("Node %d", i+1),
			weight: ws[i],
		})
	}
	return &WRRBalancer{
		nodes: nodes,
		total: tot,
	}
}

func (b *WRRBalancer) pick() *Node {
	b.mu.Lock()
	defer b.mu.Unlock()

	var res *Node
	for _, n := range b.nodes {
		n.value = n.value + n.weight
		if res == nil || n.value > res.value {
			res = n
		}
	}
	res.value -= b.total
	return res
}

func (b *WRRBalancer) Info() {
	for _, n := range b.nodes {
		log.Printf("%s: value(%d)\n", n.name, n.value)
	}
}

func TestSmoothWRR(t *testing.T) {
	balancer := NewWRRBalancer([]int{1, 2, 3})

	// simulate n requests to handle
	n := 10
	for i := 0; i < n; i++ {
		balancer.pick().Invoke()
		balancer.Info()
	}
}
