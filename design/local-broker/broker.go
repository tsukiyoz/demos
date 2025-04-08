package localbroker

import (
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Msg struct {
	Topic   string
	Content string
}

type Broker struct {
	mu sync.RWMutex

	topics map[string][]chan Msg
}

func (b *Broker) Send(msg Msg) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var eg errgroup.Group
	for _, ch := range b.topics[msg.Topic] {
		eg.Go(func() error {
			select {
			case ch <- msg:
				return nil
			default:
				return errors.New("subscriber channel is full")
			}
		})
	}

	return eg.Wait()
}

func (b *Broker) Subscribe(topic string) (<-chan Msg, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.topics == nil {
		b.topics = make(map[string][]chan Msg)
	}

	ch := make(chan Msg)
	b.topics[topic] = append(b.topics[topic], ch)

	return ch, nil
}

func (b *Broker) Close() error {
	b.mu.Lock()
	topics := b.topics
	b.topics = nil
	b.mu.Unlock()

	for _, topic := range topics {
		for _, subCh := range topic {
			close(subCh)
		}
	}

	return nil
}
