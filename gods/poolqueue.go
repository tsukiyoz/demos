package gods

import (
	"sync/atomic"
	"unsafe"
)

// PoolDequeue is a lock-free fixed-size single-producer,
// multi-consumer queue. The single producer can both push and pop
// from the head, and consumers can pop from the tail.
//
// It has the added feature that it nils out unused slots to avoid
// unnecessary retention of objects. This is important for sync.Pool,
// but not typically a property considered in the literature.
type PoolDequeue struct {
	// headTail packs together a 32-bit head index and a 32-bit
	// tail index. Both are indexes into vals modulo len(vals) - 1.
	//
	// tail = index of oldest data in queue
	// head = index of next slot to fill
	//
	// Slots in the range [tail, head) are owned by consumers.
	// A consumer continues to own a slot outside this range until
	// it nils the slot, at which ownership passes to the producer.
	headTail atomic.Uint64

	// vals is a ring buffer of interface{} values stored in this
	// dequeue. The size of this must be a power of 2.
	//
	// vals[i].typ is nil if the slot is empty and non-nil
	// otherwise. A slot is still in use until *both* the tail
	// index has moved beyond it and typ has been set to nil. This
	// is set to nil atomically by the consumer and read
	// atomically by the producer.
	vals []eface
}

type eface struct {
	typ, val unsafe.Pointer
}

const dequeueBits = 32

type dequeueNil *struct{}

// NewPoolDequeue returns a new PoolDequeue that can hold n elements.
func NewPoolDequeue(n int) *PoolDequeue {
	return &PoolDequeue{
		vals: make([]eface, n),
	}
}

func (d *PoolDequeue) unpack(ptrs uint64) (head, tail uint32) {
	const mask = 1<<dequeueBits - 1
	head = uint32((ptrs >> dequeueBits) & mask)
	tail = uint32(ptrs & mask)
	return
}

func (d *PoolDequeue) pack(head, tail uint32) uint64 {
	return uint64(head)<<dequeueBits | uint64(tail)
}

func (d *PoolDequeue) PushHead(val any) bool {
	ptrs := d.headTail.Load()
	head, tail := d.unpack(ptrs)
	if (tail+uint32(len(d.vals)))%uint32(len(d.vals)) == head {
		// Queue is full
		return false
	}
	slot := &d.vals[head&uint32(len(d.vals)-1)]

	typ := atomic.LoadPointer(&slot.typ)
	if typ != nil {
		// Another goroutine is still cleaning up the tail, so
		// the queue is actually still full.
		return false
	}

	if val == nil {
		val = dequeueNil(nil)
	}
	*(*any)(unsafe.Pointer(slot)) = val

	// Increment head. This passes ownership of slot to popTail
	// and acts as a store barrier for writing the slot.
	d.headTail.Add(1 << dequeueBits)
	return true
}

func (d *PoolDequeue) PopHead() (any, bool) {
	var slot *eface
	for {
		ptrs := d.headTail.Load()
		head, tail := d.unpack(ptrs)
		if head == tail {
			return nil, false
		}

		head--
		ptrs2 := d.pack(head, tail)
		if d.headTail.CompareAndSwap(ptrs, ptrs2) {
			slot = &d.vals[head&uint32(len(d.vals)-1)]
			break
		}
	}
	_ = slot
	// TODO
	return nil, false
}

// UNFINISHED
