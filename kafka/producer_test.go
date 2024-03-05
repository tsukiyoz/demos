package kafka

import (
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
)

var (
	addrs          = []string{"localhost:9094"}
	testTopic      = "test_topic"
	readEventTopic = "article_read_event"
)

func TestSendReadEvent(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Partitioner = sarama.NewHashPartitioner
	producer, err := sarama.NewSyncProducer(addrs, cfg)
	assert.NoError(t, err)
	for i := 0; i < 100; i++ {
		_, _, err = producer.SendMessage(&sarama.ProducerMessage{
			Topic: readEventTopic,
			Value: sarama.StringEncoder(`{"aid": 1, "uid": 123}`),
		})
	}
	assert.NoError(t, err)
}

func TestSyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Partitioner = sarama.NewHashPartitioner
	cfg.Producer.Partitioner = sarama.NewCustomPartitioner()
	producer, err := sarama.NewSyncProducer(addrs, cfg)
	assert.NoError(t, err)
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: testTopic,
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("trace_id"),
				Value: []byte("123456"),
			},
		},
		Metadata: "this is a metadata",
		Value:    sarama.StringEncoder("hello, this is tsukiyo sync speaking"),
	})
	assert.NoError(t, err)
}

func TestAsyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	producer, err := sarama.NewAsyncProducer(addrs, cfg)
	assert.NoError(t, err)
	msgCh := producer.Input()
	msgCh <- &sarama.ProducerMessage{
		Topic: testTopic,
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("trace_id"),
				Value: []byte("123456"),
			},
		},
		Metadata: "this is a metadata",
		Value:    sarama.StringEncoder("hello, this is tsukiyo async speaking"),
	}
	errCh := producer.Errors()
	succCh := producer.Successes()
	for {
		select {
		case err := <-errCh:
			t.Logf("[%v]send msg (%v) failed, err: %v\n", time.Now().Format(time.DateTime), err.Msg.Value, err.Error())
		case msg := <-succCh:
			bs, _ := msg.Value.Encode()
			t.Logf("[%v]send msg (%v) success\n", time.Now().Format(time.DateTime), string(bs))
		default:
			t.Logf("waiting for recv msg...\n")
			time.Sleep(time.Second)
		}
	}
}
