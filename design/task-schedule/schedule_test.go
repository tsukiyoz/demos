package taskschedule

import (
	"fmt"
	"testing"
	"time"
)

func TestScheduler_Run(t *testing.T) {
	scheduler := NewScheduler(stateHandler)
	defer scheduler.Close()

	printTaskStatus := func(task *Task) error {
		fmt.Printf("task %d status: %v\n", task.ID, task.Status)
		return nil
	}

	// producer
	go func() {
		ticker := time.NewTicker(20 * time.Millisecond)
		id := 6
		for {
			<-ticker.C
			scheduler.AddTask(&Task{ID: id, Status: StatusPending, Payload: map[string]string{}})
			id++
		}
	}()

	tasks := []*Task{
		{ID: 1, Status: StatusPending, Payload: map[string]string{}, Callback: printTaskStatus},
		{ID: 2, Status: StatusPending, Payload: map[string]string{}},
		{ID: 3, Status: StatusPending, Payload: map[string]string{}},
		{ID: 4, Status: StatusPending, Payload: map[string]string{}},
		{ID: 5, Status: StatusPending, Payload: map[string]string{}},
	}
	for _, task := range tasks {
		scheduler.AddTask(task) // consumer
	}

	select {}
}
