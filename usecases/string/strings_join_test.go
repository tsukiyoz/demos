package string

import (
	"strings"
	"testing"
	"unsafe"
)

var input = []string{
	"Hello World",
	"Hello  World",
	"Hello   World",
	"Hello    World",
	"Hello     World",
	"Hello      World",
	"Hello       World",
	"Hello        World",
	"Hello         World",
	"Hello          World",
}

func ConcatString(param ...string) string {
	switch len(param) {
	case 0:
		return ""
	case 1:
		return param[0]
	}
	var length int
	for _, s := range param {
		length += len(s)
	}
	bs := make([]byte, length)
	var i int
	for _, value := range param {
		i += copy(bs[i:], value)
	}
	return Bytes2Str(bs[:])
}

func Bytes2Str(slice []byte) string {
	if len(slice) < 1 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&slice))
}

func BenchmarkStringsJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = strings.Join(input, "")
	}
}

func BenchmarkUnsafe(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ConcatString(input...)
	}
}
