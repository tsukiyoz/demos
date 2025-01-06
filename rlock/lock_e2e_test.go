package rlock

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ClientE2ETestSuite struct {
	suite.Suite
	rds *redis.Pool
}

func (s *ClientE2ETestSuite) SetupSuite() {
	s.rds = &redis.Pool{
		MaxIdle:   3,
		MaxActive: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}

	res, err := s.rds.Get().Do("PING")
	s.Require().NoError(err)
	s.Require().Equal("PONG", res)
}

func TestClientE2E(t *testing.T) {
	suite.Run(t, new(ClientE2ETestSuite))
}

func (s *ClientE2ETestSuite) TestLock() {
	t := s.T()
	rds := s.rds
	tests := []struct {
		name string

		key        string
		expiration time.Duration
		retry      RetryStrategy
		timeout    time.Duration

		client *Client

		wantLock *Lock
		wantErr  error
		before   func()
		after    func()
	}{
		{
			name:       "lock success with one try",
			key:        "locked-key",
			expiration: 15 * time.Second,
			retry: &FixIntervalRetry{
				Interval: 1 * time.Second,
				MaxTries: 1,
			},
			timeout: time.Second,
			client:  NewClient(rds),
			after: func() {
				conn, err := rds.GetContext(context.Background())
				require.NoError(t, err)
				res, err := conn.Do("DEL", "locked-key")
				require.NoError(t, err)
				require.Equal(t, int64(1), res)
			},
		},
		{
			name: "lock failed",
			key:  "failed-key",
			retry: &FixIntervalRetry{
				Interval: 1 * time.Second,
				MaxTries: 3,
			},
			timeout: time.Second,
			client:  NewClient(rds),
			wantErr: ErrFailedToPreeptLock,
			before: func() {
				conn, err := rds.GetContext(context.Background())
				require.NoError(t, err)
				_, err = conn.Do("SET", "failed-key", "locked", "EX", 60)
				require.NoError(t, err)
			},
			after: func() {
				conn, err := rds.GetContext(context.Background())
				require.NoError(t, err)
				res, err := conn.Do("DEL", "failed-key")
				require.NoError(t, err)
				require.Equal(t, int64(1), res)
			},
		},
		{
			name: "already locked and lock again will refresh its expiration",
			key:  "already-locked-key",
			retry: &FixIntervalRetry{
				Interval: 1 * time.Second,
				MaxTries: 3,
			},
			timeout:    time.Second,
			expiration: 15 * time.Second,
			client:     NewClient(rds, WithValueGenFunc(func() string { return "already-locked-value" })),
			before: func() {
				conn, err := rds.GetContext(context.Background())
				require.NoError(t, err)
				_, err = conn.Do("SET", "already-locked-key", "already-locked-value", "EX", 60)
				require.NoError(t, err)
			},
			after: func() {
				conn, err := rds.GetContext(context.Background())
				require.NoError(t, err)
				res, err := conn.Do("DEL", "already-locked-key")
				require.NoError(t, err)
				require.Equal(t, int64(1), res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before()
			}

			lock, err := tt.client.Lock(context.Background(), tt.key, tt.expiration, tt.retry, tt.timeout)
			assert.True(t, errors.Is(err, tt.wantErr))
			if err != nil {
				assert.Nil(t, lock)
				return
			}

			assert.Equal(t, tt.key, lock.key)
			assert.Equal(t, tt.expiration, lock.expiration)
			assert.NotEmpty(t, lock.value)

			if tt.after != nil {
				tt.after()
			}
		})
	}
}

func (s *ClientE2ETestSuite) TestTryLock() {
	t := s.T()
	rds := s.rds
	client := NewClient(rds)

	tests := []struct {
		name string

		key        string
		expiration time.Duration

		wantLock *Lock
		wantErr  error

		before func()
		after  func()
	}{
		{
			name:       "try lock success",
			key:        "try-lock-key",
			expiration: 15 * time.Second,
			after: func() {
				conn, err := rds.GetContext(context.Background())
				require.NoError(t, err)
				res, err := conn.Do("DEL", "try-lock-key")
				require.NoError(t, err)
				require.Equal(t, int64(1), res)
			},
		},
		{
			name: "try lock failed",
			key:  "failed-key",
			before: func() {
				conn, err := rds.GetContext(context.Background())
				require.NoError(t, err)
				_, err = conn.Do("SET", "failed-key", "locked", "EX", 60)
				require.NoError(t, err)
			},
			wantErr: ErrFailedToPreeptLock,
			after: func() {
				conn, err := rds.GetContext(context.Background())
				require.NoError(t, err)
				res, err := conn.Do("DEL", "failed-key")
				require.NoError(t, err)
				require.Equal(t, int64(1), res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before()
			}

			lock, err := client.TryLock(context.Background(), tt.key, tt.expiration)
			assert.True(t, errors.Is(err, tt.wantErr))
			if err != nil {
				assert.Nil(t, lock)
				return
			}

			assert.Equal(t, tt.key, lock.key)
			assert.Equal(t, tt.expiration, lock.expiration)
			assert.NotEmpty(t, lock.value)

			if tt.after != nil {
				tt.after()
			}
		})
	}
}

func (s *ClientE2ETestSuite) TestUnlock() {
	t := s.T()
	rds := s.rds
	client := NewClient(rds)

	tests := []struct {
		name string

		lock *Lock

		before func()
		after  func()

		wantErr error
	}{
		{
			name: "unlock success",
			lock: func() *Lock {
				lock, err := client.TryLock(context.Background(), "unlocked-key", time.Minute)
				require.NoError(t, err)
				return lock
			}(),
			after: func() {
				conn, err := rds.GetContext(context.Background())
				require.NoError(t, err)
				res, err := conn.Do("EXISTS", "unlocked-key")
				require.NoError(t, err)
				require.Equal(t, int64(0), res)
			},
		},
		{
			name:    "lock not hold by the client",
			lock:    newLock(rds, "lock-not-hold-key", "lock-not-hold-value", time.Minute),
			wantErr: ErrLockNotHold,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before()
			}
			err := tt.lock.Unlock(context.Background())
			require.Equal(t, tt.wantErr, err)
			if tt.after != nil {
				tt.after()
			}
		})
	}
}
