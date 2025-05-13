package goprotobuf

import (
	"fmt"
	"os"
	"testing"

	"demo/gogoprotobuf-fast/submit"

	"github.com/gogo/protobuf/proto"
)

var request = submit.Request{
	Recvtime: 170123456,
	Uniqueid: "a1b2c3d4e5f6g7h8i9",
	Token:    "xxxx-1111-yyyy-2222-zzzz-3333",
	Phone:    "13900010002",
	Content:  "Customizing the fields of the messages to be the fields that you actually want to useremoves the need to copy between the structs you use and structs you use to serialize. gogoprotobufalso offers more serialization formats and generation of tests and even more methods.",
	Sign:     "tsukiyoXZYDFDS",
	Type:     "submit",
	Extend:   "extend",
	Version:  "v1.0.0",
}

var bs []byte

func init() {
	var err error
	bs, err = proto.Marshal(&request)
	if err != nil {
		fmt.Printf("marshal err:%s\n", err)
		os.Exit(1)
	}
}

func BenchmarkMarshal(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = proto.Marshal(&request)
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var req submit.Request
		_ = proto.Unmarshal(bs, &req)
	}
}

func BenchmarkMarshalInParallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = proto.Marshal(&request)
		}
	})
}

func BenchmarkUnmarshalInParallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var req submit.Request
			_ = proto.Unmarshal(bs, &req)
		}
	})
}
