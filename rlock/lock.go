package rlock

import (
	"context"
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	luaLock    string
	luaUnlock  string
	luaRefresh string

	ErrFailedToPreeptLock = errors.New("failed to preempt lock")
	ErrLockNotHold        = errors.New("lock not hold")
)

type Client struct {
	rds    *redis.Pool
	valuer func() string
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) exec(ctx context.Context, cmd string, args ...interface{}) (interface{}, error) {
	conn, err := c.rds.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.Do(cmd, args...)
}

func (c *Client) script(ctx context.Context, script string, args ...interface{}) (interface{}, error) {
	conn := c.rds.Get()
	if err := conn.Err(); err != nil {
		return nil, err
	}
	defer conn.Close()

	luaScript := redis.NewScript(1, script)
	return luaScript.DoContext(ctx, conn, args...)
}

func (c *Client) Lock(ctx context.Context, key string, expiration time.Duration, retry RetryStrategy, timeout time.Duration) (*Lock, error) {
	val := c.valuer()
	var timer *time.Timer
	defer func() {
		if timer != nil {
			timer.Stop()
		}
	}()
	for {
		lctx, cancel := context.WithTimeout(ctx, timeout)
		res, err := c.script(lctx, luaLock, key, val, expiration.Seconds())
		cancel()
		if err != nil && !errors.Is(err, redis.ErrNil) {
			return nil, err
		}
		if res == "OK" {
			return newLock(), nil
		}
		interval, ok := retry.Next()
		if !ok {
			if err != nil {
				err = errors.New("last retry failed: " + err.Error())
			} else {
				err = ErrLockNotHold
			}
			return nil, err
		}
		if timer == nil {
			timer = time.NewTimer(interval)
		} else {
			timer.Reset(interval)
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timer.C:
		}
	}
}

func (c *Client) TryLock(ctx context.Context, key string, expiration time.Duration) (*Lock, error) {
	val := c.valuer()
	ok, err := c.exec(ctx, "SET", key, val, "EX", int(expiration.Seconds()), "NX")
	if err != nil {
		return nil, err
	}
	if ok != "OK" {
		return nil, ErrFailedToPreeptLock
	}
	return newLock(), nil
}

type Lock struct{}

func newLock() *Lock {
	return &Lock{}
}

func (l *Lock) AutoRefresh(interval time.Duration, timeout time.Duration) error {
	// TODO
	return nil
}

func (l *Lock) Refresh(ctx context.Context) error {
	// TODO
	return nil
}

func (l *Lock) Unlock(ctx context.Context) error {
	// TODO
	return nil
}
