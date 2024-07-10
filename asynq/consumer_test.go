package main

import (
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/suite"
	"github.com/tsukaychan/demos/asynq/tasks"
	"testing"
)

func TestConsumer(t *testing.T) {
	suite.Run(t, new(ConsumerTestSuite))
}

type ConsumerTestSuite struct {
	suite.Suite
	server *asynq.Server
}

func (s *ConsumerTestSuite) SetupSuite() {
	s.T().Logf("initialing asynq server...\n")
	s.server = asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: "127.0.0.1:6379",
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)
}

func (s *ConsumerTestSuite) TestConsumer() {
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeEmailDelivery, tasks.HandleEmailDeliveryTask)
	mux.Handle(tasks.TypeImageResize, tasks.NewImageProcessor())

	if err := s.server.Run(mux); err != nil {
		s.T().Fatalf("could not run server: %v", err)
	}
}
