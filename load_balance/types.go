package test

import "log"

type Node struct {
	name          string
	weight, value int
}

func (n *Node) Invoke() {
	log.Printf("[%s], initiated a call", n.name)
}
