package goqueue

import (
	"context"
	"fmt"
	"github.com/golang-queue/queue"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestQueueUsage(t *testing.T) {
	suite.Run(t, new(QueueUsageTestSuite))
}

type QueueUsageTestSuite struct {
	suite.Suite
}

func (s *QueueUsageTestSuite) TestBasicUsageOfPool() {
	taskN := 100
	rets := make(chan string, taskN)

	q := queue.NewPool(5)
	defer q.Release()

	for i := 0; i < taskN; i++ {
		go func(i int) {
			if err := q.QueueTask(func(ctx context.Context) error {
				rets <- fmt.Sprintf("Hi Gopher, handle the job: %02d", +i)
				return nil
			}); err != nil {
				panic(err)
			}
		}(i)
	}

	for i := 0; i < taskN; i++ {
		fmt.Println("message:", <-rets)
		time.Sleep(20 * time.Millisecond)
	}
}
