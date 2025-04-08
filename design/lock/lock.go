package lock

import "sync/atomic"

const (
	UNLOCKED int32 = iota
	LOCKED
)

type Lock struct {
	state atomic.Int32
}

func (l *Lock) Lock() {
	// fast path:
	var try int32
	locked := false
	for locked = l.state.CompareAndSwap(UNLOCKED, LOCKED); !locked && try < 10; try++ {
	}
	if locked {
		return
	}

	// If we reach here, it means we failed to acquire the lock
	// runtime enqueue the goroutine to try again

	// slow path:

	// someone rewaken the goroutine here
	// and we can try again

	// there have two cas lock mode:
	// starving mode, give the lock to the starving goroutine
	// normal fair competition mode, use cas to fetch the lock
}
