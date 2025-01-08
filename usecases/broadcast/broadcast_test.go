package broadcast

import (
	"log"
	"testing"

	. "github.com/dustin/go-broadcast"
)

func TestBroadcastUsage(t *testing.T) {
	b := NewBroadcaster(100)

	workerOne(b)
}

func workerOne(b Broadcaster) {
	ch := make(chan interface{})
	b.Register(ch)
	defer b.Unregister(ch)

	// Dump out each message sent to the broadcaster.
	go func() {
		for v := range ch {
			log.Printf("workerOne read %v", v)
		}
	}()
}

func workerTwo(b Broadcaster) {
	ch := make(chan interface{})
	b.Register(ch)
	defer b.Unregister(ch)
	defer log.Printf("workerTwo is done\n")

	go func() {
		log.Printf("workerTwo read %v\n", <-ch)
	}()
}
