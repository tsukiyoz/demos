package localbroker

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestBroker_Send(t *testing.T) {
	b := &Broker{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(4)

	// simulate a producer
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				t.Log("producer done")
				b.Close()
				return
			default:
				err := b.Send(Msg{Topic: "topic", Content: time.Now().Format(time.DateTime)})
				if err != nil {
					t.Log(err)
					return
				}
				time.Sleep(time.Millisecond * 200)
			}
		}
	}()

	// simulate a consumer
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			sub, err := b.Subscribe("topic")
			if err != nil {
				t.Log(err)
				return
			}
			for msg := range sub {
				t.Log(msg.Content)
			}
			t.Log("sub done")
		}()
	}

	wg.Wait()
}
