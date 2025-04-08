package common

import (
	"errors"
	"sync/atomic"

	"golang.org/x/sys/cpu"
)

// RateLimiter is a simple rate limiter that allows a certain number of requests per second.
// It uses atomic operations to ensure thread safety and avoid contention.
// The rate limiter is designed to be used in high-performance scenarios where low latency is critical.
// It uses a 64-bit integer to store the state, where the upper 40 bits represent the timestamp
// and the lower 24 bits represent the number of requests made in the current second.
// The rate limiter allows up to 0xFFFFFF requests per second (approximately 16 million).
// The timestamp is stored in the upper 40 bits of the state, which allows for a maximum of
// 2^40 seconds (approximately 1.1 trillion years) before the timestamp wraps around.
type RateLimiter struct {
	_     cpu.CacheLinePad
	state uint64
	_     cpu.CacheLinePad
	qps   uint64
}

func NewRateLimiter(qps uint64) (*RateLimiter, error) {
	limiter := &RateLimiter{}
	if err := limiter.SetQps(qps); err != nil {
		return nil, err
	}
	return limiter, nil
}

func (limiter *RateLimiter) SetQps(qps uint64) error {
	if qps > 0xFFFFFF {
		return errors.New("rate limiter qps out of range")
	}

	atomic.StoreUint64(&limiter.qps, qps)
	return nil
}

func (limiter *RateLimiter) Take(now uint64) bool {
	for {
		state := atomic.LoadUint64(&limiter.state)
		ts := state >> 24 // 右移24位,低24位表示的是qps
		if ts < now {
			newState := now<<24 + 1 // 左移24位
			if atomic.CompareAndSwapUint64(&limiter.state, state, newState) {
				return true
			}
		} else if ts == now {
			if state&0xFFFFFF < atomic.LoadUint64(&limiter.qps) {
				newState := state + 1
				if atomic.CompareAndSwapUint64(&limiter.state, state, newState) {
					return true
				}
			} else {
				return false
			}
		} else {
			return false
		}
	}
}
