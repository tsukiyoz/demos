package taskschedule

import (
	"fmt"
	"time"
)

type TaskHandler func(*Task) error

var stateHandler = map[TaskStatus]TaskHandler{
	StatusPending:           handleDownload,
	StatusDownloadCompleted: handleAnalysis,
	StatusAnalysisCompleted: handleFinish,
}

func handleDownload(task *Task) error {
	time.Sleep(500 * time.Millisecond)
	task.Payload["image"] = "image.jpg"
	fmt.Printf("Task %d downloaded successfully.\n", task.ID)
	task.Status, _ = task.Status.Next()
	return nil
}

func handleAnalysis(task *Task) error {
	time.Sleep(1 * time.Second)
	fmt.Printf("Task %d analyzed successfully.\n", task.ID)
	task.Payload["result"] = "result.json"
	delete(task.Payload, "image")
	task.Status, _ = task.Status.Next()
	return nil
}

func handleFinish(task *Task) error {
	time.Sleep(20 * time.Millisecond)
	fmt.Printf("Task %d finished successfully.\n", task.ID)
	task.Status, _ = task.Status.Next()
	return nil
}
