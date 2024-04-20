package nats

import (
	"sync"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
)

func TestSubscriber(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	// connect to server
	nc, err := nats.Connect(nats.DefaultURL)
	require.NoError(t, err)

	defer nc.Close()

	// Simple Async Subscriber
	nc.Subscribe("foo", func(msg *nats.Msg) {
		t.Logf("Received a message: %s\n", string(msg.Data))
		wg.Done()
	})

	wg.Wait()
}

func TestSubscriberEncoded(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	// connect to server
	nc, err := nats.Connect(nats.DefaultURL)
	require.NoError(t, err)
	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	require.NoError(t, err)

	defer c.Close()

	type person struct {
		Name    string
		Address string
		Age     int
	}
	// Simple Async Subscriber
	c.Subscribe("foo", func(p *person) {
		t.Logf("Received a person: %+v\n", p)
		wg.Done()
	})

	wg.Wait()
}

func TestSubscriberResponse(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	// connect to server
	nc, err := nats.Connect(nats.DefaultURL)
	require.NoError(t, err)
	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	require.NoError(t, err)

	defer c.Close()

	c.Subscribe("help", func(subj, reply string, msg string) {
		c.Publish(reply, "i can help!")
		wg.Done()
	})

	wg.Wait()
}
