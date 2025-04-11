package errors

import (
	"errors"
	"fmt"
	"testing"
)

// 测试代码片段
func BenchmarkErrorsNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = errors.New("error message")
	}
}

func BenchmarkFmtErrorf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Errorf("error %s", "message")
	}
}

// goos: darwin
// goarch: arm64
// pkg: github.com/tsukiyoz/demos/usecases/errors
// cpu: Apple M2
// BenchmarkErrorsNew-8    1000000000               0.2937 ns/op
// BenchmarkFmtErrorf-8    24257084                49.14 ns/op
// PASS
// ok      github.com/tsukiyoz/demos/usecases/errors       2.006s
