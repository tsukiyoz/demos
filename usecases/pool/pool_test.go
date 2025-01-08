package pool

import (
	"encoding/json"
	"sync"
	"testing"
)

type User struct {
	Name   string
	Age    int32
	Remark [1024]byte
}

var buf, _ = json.Marshal(&User{
	Name: "lazywoo",
	Age:  21,
})

func TestMarshal(t *testing.T) {
	t.Logf("%s", string(buf))
}

func unmarshal() error {
	u := &User{}
	return json.Unmarshal(buf, u)
}

func BenchmarkUnmarshal(b *testing.B) {
	for n := 0; n < b.N; n++ {
		u := &User{}
		json.Unmarshal(buf, u)
	}
}

var Pool = sync.Pool{
	New: func() any {
		return new(User)
	},
}

func BenchmarkUnmarshalWithPool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		u := Pool.Get().(*User)
		json.Unmarshal(buf, u)
		Pool.Put(u)
	}
}
