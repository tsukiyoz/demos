package nsq

import (
	"net"
	"testing"
	"time"

	"github.com/nsqio/go-nsq"
)

func TestProducer(t *testing.T) {
	config := nsq.NewConfig()
	config.LocalAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:0")

	w, _ := nsq.NewProducer("127.0.0.1:7760", config)
	tk := time.NewTicker(time.Second * 3)
	for range tk.C {
		err := w.Publish("test", []byte("hello world"))
		if err != nil {
			t.Fatal(err)
		}
	}

	w.Stop()
}
