package multibytes_test

import (
	"testing"

	"github.com/tsukiyoz/demos/design/multibytes"
)

func TestMultiBytes(t *testing.T) {
	mb := multibytes.NewMultiBytes([][]byte{
		[]byte("hello, world!\n"),
	})
	_, _ = mb.Write([]byte("hello, golang!\n"))
	p := make([]byte, 32)
	_, _ = mb.Read(p)
	t.Logf("%s", p)
}
