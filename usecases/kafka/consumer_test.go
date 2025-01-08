package kafka

import (
	"context"
	"log"
	"testing"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
)

func TestConsumer(t *testing.T) {
	cfg := sarama.NewConfig()
	consumer, err := sarama.NewConsumerGroup(addrs, "test_consumer", cfg)
	assert.NoError(t, err)

	err = consumer.Consume(context.Background(), []string{testTopic}, &testConsumerHandler{})
	t.Log(err)
}

type testConsumerHandler struct{}

func (t *testConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	log.Println("setup...")
	return nil
}

func (t *testConsumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Println("cleanup...")
	return nil
}

func (t *testConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgsCh := claim.Messages()

	for msg := range msgsCh {
		log.Println(string(msg.Value))
		session.MarkMessage(msg, "")
	}

	return nil
}
