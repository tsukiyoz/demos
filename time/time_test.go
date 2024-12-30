package time_test

import (
	"testing"
	"time"
)

func TestBasicUsage(t *testing.T) {
	duration := time.Minute * 3
	t.Log(duration.Minutes())
	t.Log(duration.Seconds())
}
