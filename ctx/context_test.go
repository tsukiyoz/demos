package ctx

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// preferred usage
type CtxKey struct{}

func TestContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), "key1", "value1")
	val := ctx.Value("key1") // value1
	t.Logf("%v\n", val)

	var ctxKey CtxKey

	subCtx := context.WithValue(ctx, ctxKey, "Value")
	subCtx, cancel := context.WithTimeout(subCtx, time.Minute*10)
	t.Logf("%v\n", subCtx.Err())
	cancel()
	t.Logf("%v\n", subCtx.Err())
	t.Logf("%v\n", subCtx.Done())
	t.Logf("%v\n", ctx.Value(ctxKey))    // <nil>
	t.Logf("%v\n", subCtx.Value(ctxKey)) // Value

	subCtx = context.WithValue(subCtx, ctxKey, "ValueDifferent")
	t.Logf("%v\n", subCtx.Value(ctxKey)) // ValueDifferent
}

func TestCancel(t *testing.T) {
	// ctx := context.WithTimeout(context.Background(), time.Second * 3)
	// ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	ctx, cancel := context.WithCancel(context.Background())

	// work listen quit signal and stop loop
	go func(gCtx context.Context) {
		i := 0
		for {
			select {
			case <-ctx.Done():
				t.Logf("finish loop...")
				return
			default:
				t.Logf("%v loop at %v\n", i, time.Now().Format(time.DateTime))
				i += 1
				time.Sleep(time.Second)
			}
		}
	}(ctx)

	// control context to quit
	time.Sleep(3 * time.Second)
	cancel()
}

func TestContextErr(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	time.Sleep(time.Second * 4)
	cancel()

	switch err := ctx.Err(); {
	case errors.Is(err, context.Canceled):
		t.Logf("cancelled!\n")
	case errors.Is(err, context.DeadlineExceeded):
		t.Logf("dealine exceeded!\n")
	}
}

func TestContextFatherCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	subCtx, _ := context.WithCancel(ctx)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 3)
		t.Logf("father ctx cancelled!\n")
		cancel()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			switch err := subCtx.Err(); {
			case errors.Is(err, context.Canceled):
				t.Logf("son ctx has been cancelled!\n")
				return
			case errors.Is(err, context.DeadlineExceeded):
				t.Logf("son context not cancelled until timeout!\n")
				return
			default:
				t.Logf("son ctx waiting for cancel...\n")
				time.Sleep(time.Second)
			}
		}
	}()

	wg.Wait()
}

func TestContextSonCancel(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*7)
	_, sonCancel := context.WithCancel(ctx)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		time.Sleep(time.Second * 3)
		t.Logf("son ctx cancelled!\n")
		sonCancel()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			switch err := ctx.Err(); {
			case errors.Is(err, context.Canceled):
				t.Logf("father context has been cancelled!\n")
				return
			case errors.Is(err, context.DeadlineExceeded):
				t.Logf("father context not cancelled until timeout!\n")
				return
			default:
				time.Sleep(time.Second)
				t.Logf("father ctx waiting for cancel...\n")
			}
		}
	}()

	wg.Wait()
}
