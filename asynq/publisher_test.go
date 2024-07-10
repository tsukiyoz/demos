package main

import (
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/suite"
	"github.com/tsukaychan/demos/asynq/tasks"
	"testing"
	"time"
)

func TestPublisher(t *testing.T) {
	suite.Run(t, new(PublisherTestSuite))
}

type PublisherTestSuite struct {
	suite.Suite
	client *asynq.Client
}

func (s *PublisherTestSuite) SetupSuite() {
	s.T().Logf("initialing asynq client...\n")
	s.client = asynq.NewClient(asynq.RedisClientOpt{
		Addr: "127.0.0.1:6379",
	})
}

func (s *PublisherTestSuite) TearDownTest() {
	s.T().Logf("teardown asynq client...\n")
	s.client.Close()
}

func (s *PublisherTestSuite) TestPublishTask() {
	task, err := tasks.NewEmailDeliveryTask(42, "some:template:id")
	if err != nil {
		s.T().Fatalf("could not create task: %v", err)
	}

	info, err := s.client.Enqueue(task)
	if err != nil {
		s.T().Fatalf("cound not enqueue task: %v", err)
	}

	s.T().Logf("enqueued task: id=%s queue=%s\n", info.ID, info.Queue)
}

func (s *PublisherTestSuite) TestPublishDelayedTask() {
	task, err := tasks.NewEmailDeliveryTask(42, "some:template:id")
	if err != nil {
		s.T().Fatalf("could not create task: %v", err)
	}

	info, err := s.client.Enqueue(task, asynq.ProcessIn(time.Second*5))
	if err != nil {
		s.T().Fatalf("cound not enqueue task: %v", err)
	}

	s.T().Logf("enqueued task: id=%s queue=%s\n", info.ID, info.Queue)
}

func (s *PublisherTestSuite) TestBatchPublishTask() {
	n := 500
	for i := 0; i < n; i++ {
		task, err := tasks.NewEmailDeliveryTask(42, "some:template:id")
		if err != nil {
			s.T().Fatalf("could not create task: %v", err)
		}
		info, err := s.client.Enqueue(task, asynq.ProcessIn(time.Second*5), asynq.MaxRetry(3))
		if err != nil {
			s.T().Fatalf("cound not enqueue task: %v", err)
		}
		s.T().Logf("enqueued task: id=%s queue=%s\n", info.ID, info.Queue)
	}
}
