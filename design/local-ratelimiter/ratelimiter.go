package common

import (
	"errors"
	"sync/atomic"

	"golang.org/x/sys/cpu"
)

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
