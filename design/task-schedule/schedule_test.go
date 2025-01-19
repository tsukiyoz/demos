package taskschedule

import (
	"fmt"
	"testing"
	"time"
)

// ----------- task ------------

type AnalysisTask struct {
	id      int
	status  AnalysisTaskStatus
	Payload map[string]string
}

func (t *AnalysisTask) ID() int {
	return t.id
}

func (t *AnalysisTask) Status() TaskStatus {
	return t.status.TaskStatus()
}

func (t *AnalysisTask) Next() bool {
	switch t.status {
	case StatusPending:
		return true
	case StatusDownloadCompleted:
		return true
	case StatusAnalysisCompleted:
		return true
	case StatusAllCompleted:
		return false
	default:
		return false
	}
}

func (t *AnalysisTask) Callback() func(Task) error {
	return nil
}

type AnalysisTaskStatus int

const (
	StatusUnknown AnalysisTaskStatus = iota
	StatusPending
	StatusDownloadCompleted
	StatusAnalysisCompleted
	StatusAllCompleted
)

func (s AnalysisTaskStatus) TaskStatus() TaskStatus {
	return TaskStatus(s)
}

// ----------- handler ------------

var stateHandler = map[TaskStatus]TaskHandler[*AnalysisTask]{
	StatusPending.TaskStatus():           handleDownload,
	StatusDownloadCompleted.TaskStatus(): handleAnalysis,
	StatusAnalysisCompleted.TaskStatus(): handleFinish,
}

func handleDownload(task *AnalysisTask) error {
	time.Sleep(500 * time.Millisecond)
	task.Payload["image"] = "image.jpg"
	fmt.Printf("Task %d downloaded successfully.\n", task.ID())
	task.status = StatusDownloadCompleted
	return nil
}

func handleAnalysis(task *AnalysisTask) error {
	time.Sleep(1 * time.Second)
	fmt.Printf("Task %d analyzed successfully.\n", task.ID())
	task.Payload["result"] = "result.json"
	delete(task.Payload, "image")
	task.status = StatusAnalysisCompleted
	return nil
}

func handleFinish(task *AnalysisTask) error {
	time.Sleep(20 * time.Millisecond)
	fmt.Printf("Task %d finished successfully.\n", task.ID())
	task.status = StatusAllCompleted
	return nil
}

// ------------- test --------------

func TestScheduler_Run(t *testing.T) {
	scheduler := NewScheduler(stateHandler)
	defer scheduler.Close()

	// producer
	go func() {
		ticker := time.NewTicker(20 * time.Millisecond)
		id := 6
		for {
			<-ticker.C
			scheduler.AddTask(&AnalysisTask{id: id, status: StatusPending, Payload: map[string]string{}})
			id++
		}
	}()

	tasks := []*AnalysisTask{
		{id: 1, status: StatusPending, Payload: map[string]string{}},
		{id: 2, status: StatusPending, Payload: map[string]string{}},
		{id: 3, status: StatusPending, Payload: map[string]string{}},
		{id: 4, status: StatusPending, Payload: map[string]string{}},
		{id: 5, status: StatusPending, Payload: map[string]string{}},
	}
	for _, task := range tasks {
		scheduler.AddTask(task) // consumer
	}

	select {}
}
