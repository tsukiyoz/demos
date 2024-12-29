package main

import (
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/tsukaychan/demos/asynq/tasks"
)

func startConsumer() error {
	server := asynq.NewServer(
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
			IsFailure: func(err error) bool {
				return err == tasks.ErrProcessFailed
			},
			RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
				return asynq.DefaultRetryDelayFunc(n, e, t)
			},
		},
	)
	defer server.Shutdown()

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeEmailDelivery, tasks.HandleEmailDeliveryTask)
	mux.Handle(tasks.TypeImageResize, tasks.NewImageProcessor())

	if err := server.Run(mux); err != nil {
		fmt.Printf("could not run server: %v", err)
	}

	return nil
}
