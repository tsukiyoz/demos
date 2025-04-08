package bufferpool

import (
	"slices"
	"sync"
	"sync/atomic"
)

const (
	minBitSize = 6 // 2**6=64 is a CPU cache line size
	steps      = 20

	minSize = 1 << minBitSize
	maxSize = 1 << (minBitSize + steps - 1)

	calibrateCallsThreshold = 42000
	maxPercentile           = 0.95
)

type Pool struct {
	calls       [steps]uint64
	calibrating uint64

	defaultSize uint64
	maxSize     uint64

	pool sync.Pool
}

var defaultPool Pool

func Get() *ByteBuffer {
	return defaultPool.Get()
}

func (p *Pool) Get() *ByteBuffer {
	v := p.pool.Get()
	if v == nil {
		return &ByteBuffer{
			B: make([]byte, 0, atomic.LoadUint64(&p.defaultSize)),
		}
	}
	return v.(*ByteBuffer)
}

func Put(b *ByteBuffer) {
	defaultPool.Put(b)
}

func (p *Pool) Put(b *ByteBuffer) {
	idx := index(len(b.B))

	if atomic.AddUint64(&p.calls[idx], 1) > calibrateCallsThreshold {
		p.calibrate()
	}

	maxSize := int(atomic.LoadUint64(&p.maxSize))
	if maxSize == 0 || cap(b.B) <= maxSize {
		b.Reset()
		p.pool.Put(b)
	}
}

// [2**k, 2**(k+1)-1)] => k - minBitSize
// [64, 127] => 0
// [128, 255] => 1
// [256, 511] => 2
// ...
func index(n int) int {
	n--
	n >>= minBitSize
	idx := 0
	for n > 0 {
		n >>= 1
		idx++
	}
	if idx >= steps {
		return steps - 1
	}
	return idx
}

type CallSize struct {
	calls uint64
	size  uint64
}

func (p *Pool) calibrate() {
	if !atomic.CompareAndSwapUint64(&p.calibrating, 0, 1) {
		return
	}

	a := make([]CallSize, 0, steps)
	sum := uint64(0)
	for i := uint64(0); i < steps; i++ {
		a = append(a, CallSize{
			calls: atomic.SwapUint64(&p.calls[i], 0),
			size:  minSize << i,
		})
		sum += a[i].calls
	}

	slices.SortFunc(a, func(x, y CallSize) int {
		return int(x.calls - y.calls)
	})

	defaultSize := a[0].size
	maxSize := defaultSize

	maxSum := uint64(float64(sum) * maxPercentile)
	callsSum := uint64(0)

	for i := 0; i < steps; i++ {
		if callsSum > maxSum {
			break
		}

		callsSum += a[i].calls
		maxSize = max(maxSize, a[i].size)
	}

	atomic.StoreUint64(&p.defaultSize, defaultSize)
	atomic.StoreUint64(&p.maxSize, maxSize)
	atomic.StoreUint64(&p.calibrating, 0)
}
