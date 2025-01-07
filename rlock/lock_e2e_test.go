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
		MaxIdle:   128,
		MaxActive: 256,
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
			key:        "test::lock::locked-key",
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
				res, err := conn.Do("DEL", "test::lock::locked-key")
				require.NoError(t, err)
				require.Equal(t, int64(1), res)
			},
		},
		{
			name: "lock failed",
			key:  "test::lock::failed-key",
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
				_, err = conn.Do("SET", "test::lock::failed-key", "locked", "EX", 60)
				require.NoError(t, err)
			},
			after: func() {
				conn, err := rds.GetContext(context.Background())
				require.NoError(t, err)
				res, err := conn.Do("DEL", "test::lock::failed-key")
				require.NoError(t, err)
				require.Equal(t, int64(1), res)
			},
		},
		{
			name: "already locked and lock again will refresh its expiration",
			key:  "test::lock::already-locked-key",
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
				_, err = conn.Do("SET", "test::lock::already-locked-key", "already-locked-value", "EX", 60)
				require.NoError(t, err)
			},
			after: func() {
				conn, err := rds.GetContext(context.Background())
				require.NoError(t, err)
				res, err := conn.Do("DEL", "test::lock::already-locked-key")
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
			key:        "test::trylock::try-lock-key",
			expiration: 15 * time.Second,
			after: func() {
				conn, err := rds.GetContext(context.Background())
				require.NoError(t, err)
				res, err := conn.Do("DEL", "test::trylock::try-lock-key")
				require.NoError(t, err)
				require.Equal(t, int64(1), res)
			},
		},
		{
			name: "try lock failed",
			key:  "test::trylock::failed-key",
			before: func() {
				conn, err := rds.GetContext(context.Background())
				require.NoError(t, err)
				_, err = conn.Do("SET", "test::trylock::failed-key", "locked", "EX", 60)
				require.NoError(t, err)
			},
			wantErr: ErrFailedToPreeptLock,
			after: func() {
				conn, err := rds.GetContext(context.Background())
				require.NoError(t, err)
				res, err := conn.Do("DEL", "test::trylock::failed-key")
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
				lock, err := client.TryLock(context.Background(), "test::unlock::unlocked-key", time.Minute)
				require.NoError(t, err)
				return lock
			}(),
			after: func() {
				conn, err := rds.GetContext(context.Background())
				require.NoError(t, err)
				res, err := conn.Do("EXISTS", "test::unlock::unlocked-key")
				require.NoError(t, err)
				require.Equal(t, int64(0), res)
			},
		},
		{
			name:    "lock not hold by the client",
			lock:    newLock(rds, "test::unlock::lock-not-hold-key", "test::unlock::lock-not-hold-value", time.Minute),
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

func (s *ClientE2ETestSuite) TestRefresh() {
	t := s.T()

	tests := []struct {
		name string

		timeout time.Duration
		lock    *Lock

		wantErr error

		before func()
		after  func()
	}{
		{
			name:    "refresh success",
			timeout: time.Minute,
			lock: &Lock{
				rds:        s.rds,
				key:        "test::refresh::refresh-key",
				value:      "refresh-value",
				expiration: time.Minute,
				unlock:     make(chan struct{}, 1),
			},
			before: func() {
				conn, err := s.rds.GetContext(context.Background())
				require.NoError(t, err)
				_, err = conn.Do("SET", "test::refresh::refresh-key", "refresh-value", "EX", 10)
				require.NoError(t, err)
			},
			after: func() {
				res, err := s.rds.Get().Do("TTL", "test::refresh::refresh-key")
				require.NoError(t, err)
				require.Greater(t, res.(int64), int64(50))
				_, err = s.rds.Get().Do("DEL", "test::refresh::refresh-key")
				require.NoError(t, err)
			},
		},
		{
			name:    "lock is not hold by the client",
			timeout: time.Minute,
			lock: &Lock{
				rds:        s.rds,
				key:        "test::refresh::not-hold-key",
				value:      "not-hold-value",
				expiration: time.Minute,
				unlock:     make(chan struct{}, 1),
			},
			wantErr: ErrLockNotHold,
			before: func() {
				conn, err := s.rds.GetContext(context.Background())
				require.NoError(t, err)
				_, err = conn.Do("SET", "test::refresh::not-hold-key", "another-value", "EX", 10)
				require.NoError(t, err)
			},
			after: func() {
				res, err := s.rds.Get().Do("GET", "test::refresh::not-hold-key")
				require.NoError(t, err)
				require.Equal(t, "another-value", string(res.([]byte)))
				_, err = s.rds.Get().Do("DEL", "test::refresh::not-hold-key")
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before()
			}
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			err := tt.lock.Refresh(ctx)
			cancel()
			assert.Equal(t, tt.wantErr, err)
			if tt.after != nil {
				tt.after()
			}
		})
	}
}

func (s *ClientE2ETestSuite) TestAutoRefresh() {
	t := s.T()
	client := NewClient(s.rds)

	ctx := context.Background()
	l, err := client.TryLock(ctx, "test::auto-refresh::auto-refresh-key", 15*time.Second)
	require.NoError(t, err)
	go func() {
		err = l.AutoRefresh(3*time.Second, 300*time.Millisecond)
		require.NoError(t, err)
	}()
	time.Sleep(15 * time.Second)
	res, err := s.rds.Get().Do("EXISTS", "test::auto-refresh::auto-refresh-key")
	require.NoError(t, err)
	require.Equal(t, int64(1), res)

	err = l.Unlock(ctx)
	require.NoError(t, err)
}
