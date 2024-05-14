package test

import (
	"sync/atomic"
	"testing"
)

type RRBalancer struct {
	nodes []*Node
	next  uint32
}

func NewRrBalancer(nodes []*Node) *RRBalancer {
	return &RRBalancer{
		nodes: nodes,
	}
}

func (r *RRBalancer) Pick() *Node {
	nodesLen := uint32(len(r.nodes))
	nextIdx := atomic.AddUint32(&r.next, 1)
	return r.nodes[nextIdx%nodesLen]
}

func TestRR(t *testing.T) {
	balancer := NewRrBalancer([]*Node{
		{
			name: "node1",
		},
		{
			name: "node2",
		},
		{
			name: "node3",
		},
	})

	n := 10
	for i := 0; i < n; i++ {
		balancer.Pick().Invoke()
	}
}
