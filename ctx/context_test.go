package ctx

import (
	"context"
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
	t.Logf("%v\n", ctx.Value(ctxKey))    // <nil>
	t.Logf("%v\n", subCtx.Value(ctxKey)) // Value

	subCtx = context.WithValue(subCtx, ctxKey, "ValueDifferent")
	t.Logf("%v\n", subCtx.Value(ctxKey)) // ValueDifferent
}

func TestCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
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
	time.Sleep(3 * time.Second)
	cancel()
}
