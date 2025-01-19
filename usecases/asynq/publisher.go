package main

import (
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/tsukiyoz/demos/usecases/asynq/tasks"
)

func startPublisher() error {
	opt := asynq.RedisClientOpt{
		Addr: "127.0.0.1:6379",
	}
	client := asynq.NewClient(opt)
	defer client.Close()

	// send common task
	task, err := tasks.NewEmailDeliveryTask(42, "this is a common task")
	if err != nil {
		fmt.Printf("could not create task: %v", err)
	}
	info, err := client.Enqueue(task)
	if err != nil {
		fmt.Printf("could not enqueue task: %v", err)
	}
	fmt.Printf("enqueued task: id=%s queue=%s\n", info.ID, info.Queue)

	task, err = tasks.NewImageResizeTask("https://example")
	if err != nil {
		fmt.Printf("could not create task: %v", err)
	}
	info, err = client.Enqueue(task)
	if err != nil {
		fmt.Printf("could not enqueue task: %v", err)
	}
	fmt.Printf("enqueued task: id=%s queue=%s\n", info.ID, info.Queue)

	// send delayed task
	task, err = tasks.NewEmailDeliveryTask(42, "this is a delayed task")
	if err != nil {
		fmt.Printf("could not create task: %v", err)
	}
	info, err = client.Enqueue(task, asynq.ProcessIn(5*time.Second))
	if err != nil {
		fmt.Printf("could not enqueue task: %v", err)
	}
	fmt.Printf("enqueued task: id=%s queue=%s\n", info.ID, info.Queue)

	// send cron task
	task, err = tasks.NewEmailDeliveryTask(42, "this is a cron task")
	if err != nil {
		fmt.Printf("could not create task: %v", err)
	}
	scheduler := asynq.NewScheduler(opt, nil)

	entryId, err := scheduler.Register("@every 3s", task)
	if err != nil {
		fmt.Printf("could not register task: %v", err)
	}
	fmt.Printf("registered task: id=%s\n", entryId)

	go func() {
		time.Sleep(10 * time.Second)
		err := scheduler.Unregister(entryId)
		if err != nil {
			fmt.Printf("could not unregister task: %v", err)
		}
		fmt.Printf("unregistered task: id=%s\n", entryId)
	}()

	scheduler.Run()
	defer scheduler.Shutdown()

	return nil
}
