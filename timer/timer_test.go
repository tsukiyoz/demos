package timer

import (
	"context"
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for heartbeat := range ticker.C {
		t.Logf("%v\n", heartbeat.Format(time.DateTime))
	}
}

func TestTimer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	for {
		select {
		case heartbeat := <-timer.C:
			t.Logf("%v\n", heartbeat.Format(time.DateTime))
		case <-ctx.Done():
			t.Logf("context cancelled!\n")
			return
		}
	}
}
