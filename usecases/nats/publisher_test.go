package nats

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/nats-io/nats.go"
)

func TestPublisher(t *testing.T) {
	// connect to server
	nc, err := nats.Connect(nats.DefaultURL)
	require.NoError(t, err)

	defer nc.Close()

	// Simple Publisher
	err = nc.Publish("foo", []byte("Hello World"))
	require.NoError(t, err)
}

func TestPublisherEncoded(t *testing.T) {
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
	// Simple Publisher
	err = c.Publish("foo", &person{
		Name:    "tsukiyo",
		Address: "china",
		Age:     18,
	})
	require.NoError(t, err)
}

func TestPublisherResponse(t *testing.T) {
	nc, err := nats.Connect(nats.DefaultURL)
	require.NoError(t, err)
	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	require.NoError(t, err)
	defer c.Close()

	var resp string
	err = c.Request("help", "help me", &resp, time.Second)
	require.NoError(t, err)

	t.Logf("%v\n", resp)
}
