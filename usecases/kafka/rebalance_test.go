package kafka

import (
	"context"
	"log"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/IBM/sarama"
)

// 分区状态（每个分区独立维护）
type PartitionState struct {
	mu               sync.Mutex
	processedOffsets map[int64]struct{} // 已处理的 offset 列表
	maxCommitted     int64              // 当前已提交的最大连续 offset
	wg               sync.WaitGroup     // 等待所有消息处理完成
}

// 处理消息完成时调用
func (s *PartitionState) MarkProcessed(offset int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.wg.Add(1)
	defer s.wg.Done()

	s.processedOffsets[offset] = struct{}{}

	for {
		next := s.maxCommitted + 1
		if _, ok := s.processedOffsets[next]; !ok {
			break
		}
		delete(s.processedOffsets, next)
		s.maxCommitted = next
	}
}

type ConsumerHandler struct {
	maxConcurrent int                       // 最大并发数
	semaphore     chan struct{}             // 信号量控制并发
	stateMap      map[int32]*PartitionState // 分区状态
	mu            sync.Mutex
}

func NewConsumerHandler(maxConcurrent int) *ConsumerHandler {
	return &ConsumerHandler{
		maxConcurrent: maxConcurrent,
		semaphore:     make(chan struct{}, maxConcurrent),
		stateMap:      make(map[int32]*PartitionState),
	}
}

func (h *ConsumerHandler) getState(partition int32) *PartitionState {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.stateMap[partition]; !ok {
		h.stateMap[partition] = &PartitionState{
			processedOffsets: make(map[int64]struct{}),
		}
	}
	return h.stateMap[partition]
}

func (h *ConsumerHandler) Setup(sess sarama.ConsumerGroupSession) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	claims := sess.Claims()[Topic]
	wg := sync.WaitGroup{}
	for partition, state := range h.stateMap {
		if !slices.Contains(claims, partition) {
			// rebalance后，此分区已经不属于当前消费者组
			// 可以将数据清理，避免不必要的commit
			wg.Add(1)
			go func() {
				defer wg.Done()
				state.mu.Lock()
				state.wg.Wait()              // 等待所有消息处理完成
				state.processedOffsets = nil // 清空已处理的 offset 列表
				delete(h.stateMap, partition)
				state.mu.Unlock()
			}()
		}
	}

	wg.Wait() // 等待所有分区的处理完成
	h.StartCommitLoop(sess, 3*time.Second)
	return nil
}

func (h *ConsumerHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	for partition, state := range h.stateMap {
		state.mu.Lock()
		offsetToCommit := state.maxCommitted + 1
		sess.MarkOffset(Topic, partition, offsetToCommit, "")
		state.mu.Unlock()
	}
	return nil
}

func (h *ConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		state := h.getState(msg.Partition)
		h.processMessage(msg, state)
	}
	return nil
}

func (h *ConsumerHandler) processMessage(msg *sarama.ConsumerMessage, state *PartitionState) {
	h.semaphore <- struct{}{}
	go func() {
		defer func() { <-h.semaphore }()

		// 处理消息
		// processMessage(msg)

		state.MarkProcessed(msg.Offset)
	}()
}

func (h *ConsumerHandler) StartCommitLoop(sess sarama.ConsumerGroupSession, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				snapshot := make(map[int32]int64)
				h.mu.Lock()
				for partition, state := range h.stateMap {
					state.mu.Lock()
					snapshot[partition] = state.maxCommitted
					state.mu.Unlock()
				}
				h.mu.Unlock()

				for partition, offset := range snapshot {
					sess.MarkOffset(Topic, partition, offset+1, "")
				}
			case <-sess.Context().Done():
				return
			}
		}
	}()
}

const (
	Topic = "your_topic"
)

func TestKafkaAsyncConsumer(t *testing.T) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Consumer.Offsets.AutoCommit.Enable = false // 关闭自动提交
	config.Consumer.Group.Session.Timeout = 30 * time.Second
	handler := NewConsumerHandler(100)
	group, err := sarama.NewConsumerGroup([]string{"localhost:9092"}, "my-group", config)
	if err != nil {
		log.Fatal(err)
	}
	defer group.Close()
	ctx := context.Background()
	for {
		if err := group.Consume(ctx, []string{Topic}, handler); err != nil {
			log.Fatal("消费失败:", err)
		}
	}
}
