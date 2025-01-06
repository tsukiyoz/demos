package rlock

import (
	"context"
	_ "embed"
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

const (
	success int64 = iota + 1
	failed
)

var (
	//go:embed script/lua/lock.lua
	luaLock string
	//go:embed script/lua/try_lock.lua
	luaTryLock string
	//go:embed script/lua/unlock.lua
	luaUnlock string
	//go:embed script/lua/refresh.lua
	luaRefresh string

	ErrFailedToPreeptLock = errors.New("failed to preempt lock")
	ErrLockNotHold        = errors.New("lock not hold")
)

type GenFunc func() string

type Client struct {
	rds          *redis.Pool
	valueGenFunc GenFunc
}

type Option func(*Client)

func WithValueGenFunc(getter GenFunc) Option {
	return func(c *Client) {
		c.valueGenFunc = getter
	}
}

func NewClient(rds *redis.Pool, opts ...Option) *Client {
	cli := &Client{
		rds: rds,
		valueGenFunc: func() string {
			return uuid.New().String()
		},
	}

	for _, opt := range opts {
		opt(cli)
	}

	return cli
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

func (c *Client) Lock(ctx context.Context, key string, exp time.Duration, retry RetryStrategy, timeout time.Duration) (*Lock, error) {
	val := c.valueGenFunc()
	var timer *time.Timer
	defer func() {
		if timer != nil {
			timer.Stop()
		}
	}()

	for {
		lctx, cancel := context.WithTimeout(ctx, timeout)
		secs := exp.Seconds()
		res, err := c.script(lctx, luaLock, key, val, secs)
		cancel()
		if err != nil && !errors.Is(err, redis.ErrPoolExhausted) {
			return nil, err
		}

		switch res := res.(type) {
		case int64:
			if res == success {
				return newLock(c.rds, key, val, exp), nil
			}
		default:
			panic("unexpected type")
		}

		interval, ok := retry.Next()
		if !ok {
			if err != nil {
				err = errors.New("last retry failed: " + err.Error())
			} else {
				err = ErrFailedToPreeptLock
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

func (c *Client) TryLock(ctx context.Context, key string, exp time.Duration) (*Lock, error) {
	val := c.valueGenFunc()
	res, err := c.script(ctx, luaTryLock, key, val, exp.Seconds())
	if err != nil {
		return nil, err
	}
	if res != success {
		return nil, ErrFailedToPreeptLock
	}
	return newLock(c.rds, key, val, exp), nil
}

type Lock struct {
	rds        *redis.Pool
	key        string
	value      string
	expiration time.Duration
	unlock     chan struct{}
}

func newLock(rds *redis.Pool, key string, val string, exp time.Duration) *Lock {
	return &Lock{
		rds:        rds,
		key:        key,
		value:      val,
		expiration: exp,
	}
}

func (l *Lock) script(ctx context.Context, script string, args ...interface{}) (interface{}, error) {
	conn := l.rds.Get()
	if err := conn.Err(); err != nil {
		return nil, err
	}
	defer conn.Close()

	luaScript := redis.NewScript(1, script)
	return luaScript.DoContext(ctx, conn, args...)
}

func (l *Lock) AutoRefresh(interval time.Duration, timeout time.Duration) error {
	ticker := time.NewTicker(interval)
	retryCh := make(chan struct{}, 1)
	defer func() {
		ticker.Stop()
		close(retryCh)
	}()

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			err := l.Refresh(ctx)
			cancel()
			if err == context.DeadlineExceeded {
				select {
				case retryCh <- struct{}{}:
				default:
				}
				continue
			} else {
				return err
			}
		case <-retryCh:
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			err := l.Refresh(ctx)
			cancel()
			if err == context.DeadlineExceeded {
				select {
				case retryCh <- struct{}{}:
				default:
				}
				continue
			} else {
				return err
			}
		}
	}
}

func (l *Lock) Refresh(ctx context.Context) error {
	res, err := l.script(ctx, luaRefresh, l.key, l.value)
	if err != nil {
		return err
	}
	if res != success {
		return ErrLockNotHold
	}
	return nil
}

func (l *Lock) Unlock(ctx context.Context) error {
	res, err := l.script(ctx, luaUnlock, l.key, l.value)
	if err == redis.ErrNil {
		return ErrLockNotHold
	}
	if err != nil {
		return err
	}
	if res != success {
		return ErrLockNotHold
	}
	return nil
}
