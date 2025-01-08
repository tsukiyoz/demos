package time

import (
	"testing"
	"time"
)

func TestUnixUsage(t *testing.T) {
	t.Logf(time.UnixMilli(0).Format(time.DateTime)) // 1970-01-01 08:00:00
	var null time.Time
	t.Logf(null.Format(time.DateTime)) // 0001-01-01 00:00:00
}
