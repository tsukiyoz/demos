package hierarchy

import (
	"sort"
	"time"

	. "github.com/tsukiyoz/demos/usecases/concurrency/deadlocks_train/common"
)

func lockIntersectionsInDistance(id, reserveStart, reserveEnd int, crossings []*Crossing) {
	var intersectionsToLock []*Intersection
	for _, crossing := range crossings {
		if reserveEnd >= crossing.Position && reserveStart <= crossing.Position && crossing.Intersection.LockedBy != id {
			intersectionsToLock = append(intersectionsToLock, crossing.Intersection)
		}
	}

	sort.Slice(intersectionsToLock, func(i, j int) bool {
		return intersectionsToLock[i].ID < intersectionsToLock[j].ID
	})

	for _, it := range intersectionsToLock {
		it.Mutex.Lock()
		it.LockedBy = id
		time.Sleep(10 * time.Millisecond)
	}
}

func MoveTrain(train *Train, distance int, crossings []*Crossing) {
	for train.Front < distance {
		train.Front += 1
		for _, crossing := range crossings {
			if train.Front == crossing.Position {
				lockIntersectionsInDistance(train.ID, crossing.Position, crossing.Position+train.Length, crossings)
			}
			back := train.Front - train.Length
			if back == crossing.Position {
				crossing.Intersection.LockedBy = -1
				crossing.Intersection.Mutex.Unlock()
			}
		}
		time.Sleep(30 * time.Millisecond)
	}
}
