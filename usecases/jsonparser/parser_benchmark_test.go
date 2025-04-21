package jsonparser

import (
	"testing"

	"github.com/bytedance/sonic"
	jsoniter "github.com/json-iterator/go"
	"github.com/mailru/easyjson"
)

func BenchmarkSonic_Marshal(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := sonic.Marshal(&twitter)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSonic_Unmarshal(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		err := sonic.Unmarshal(data, &twitter)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEasyJson_Marshal(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := easyjson.Marshal(twitter)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEasyJson_Unmarshal(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		err := easyjson.Unmarshal(data, twitter)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsoniter_Marshal(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		err := jsoniter.Unmarshal(data, &twitter)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsoniter_Unmarshal(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		err := jsoniter.Unmarshal(data, &twitter)
		if err != nil {
			b.Fatal(err)
		}
	}
}
