package rlock

import "time"

type RetryStrategy interface {
	Next() (time.Duration, bool)
}

type FixIntervalRetry struct {
	Interval time.Duration
	MaxTries int
	tryCnt   int
}

func (f *FixIntervalRetry) Next() (time.Duration, bool) {
	f.tryCnt++
	return f.Interval, f.tryCnt < f.MaxTries
}
