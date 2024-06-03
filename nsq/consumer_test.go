package nsq

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/nsqio/go-nsq"
)

type MessageHandler struct{}

func (h *MessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	log.Printf("Received a message: %s, at %v", string(m.Body), time.Now().Format(time.DateTime))

	return nil
}

func TestConsumer(t *testing.T) {
	cfg := nsq.NewConfig()
	cfg.LocalAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	consumer, err := nsq.NewConsumer("test", "ch", cfg)
	if err != nil {
		log.Fatal(err)
	}

	consumer.AddHandler(&MessageHandler{})

	err = consumer.ConnectToNSQD("127.0.0.1:4150")
	if err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	consumer.Stop()
}
