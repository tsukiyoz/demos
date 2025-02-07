package common

import "sync"

type Train struct {
	ID     int
	Length int
	Front  int
}

type Intersection struct {
	ID       int
	Mutex    sync.Mutex
	LockedBy int
}

type Crossing struct {
	Position     int
	Intersection *Intersection
}
