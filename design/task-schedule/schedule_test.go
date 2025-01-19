package taskschedule

import (
	"fmt"
	"testing"
	"time"
)

// ----------- task ------------

type AnalysisTask struct {
	ID      int
	status  AnalysisTaskStatus
	Payload map[string]string
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
	if time.Now().UnixMilli()%7 == 0 {
		return fmt.Errorf("download failed") // mock download failed
	}
	task.Payload["image"] = "image.jpg"
	fmt.Printf("Task %d downloaded successfully.\n", task.ID)
	task.status = StatusDownloadCompleted
	return nil
}

func handleAnalysis(task *AnalysisTask) error {
	if task.Payload["image"] == "" {
		// no image find in local memory, fetch image from oss
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("Task %d fetched image from oss.\n", task.ID)
		task.Payload["image"] = "image_from_oss.jpg"
	} else {
		fmt.Printf("Task %d found image in local memory.\n", task.ID)
	}
	time.Sleep(800 * time.Millisecond)
	fmt.Printf("Task %d analyzed successfully.\n", task.ID)
	task.Payload["result"] = "result.json"
	delete(task.Payload, "image")
	task.status = StatusAnalysisCompleted
	return nil
}

func handleFinish(task *AnalysisTask) error {
	time.Sleep(20 * time.Millisecond)
	fmt.Printf("Task %d finished successfully.\n", task.ID)
	task.status = StatusAllCompleted
	return nil
}

// ------------- test --------------

func TestScheduler_Run(t *testing.T) {
	onErrorFunc := func(task *AnalysisTask, err error) {
		fmt.Printf("Task %d failed: %v\n", task.ID, err)
	}
	var customStrategy ExecuteStrategy[*AnalysisTask]
	customStrategy = func(s *Scheduler[*AnalysisTask], task *AnalysisTask) error {
		handler, ok := stateHandler[task.Status()]
		if !ok {
			panic("unknown status")
		}
		err := handler(task)
		if err != nil {
			go func() {
				time.After(1 * time.Second)
				fmt.Printf("Task %d retrying...\n", task.ID)
				s.Submit(task)
			}()
			return err
		}
		if !task.Next() {
			return nil
		}

		if time.Now().UnixMilli()%4 != 0 {
			// greedy strategy
			customStrategy(s, task)
		} else {
			s.Submit(task)
		}

		return nil
	}

	opts := []Option[*AnalysisTask]{}
	opts = append(opts, WithOnError(onErrorFunc))
	// opts = append(opts, WithExecuteStrategy(NewParallelStrategy[*AnalysisTask]()))
	opts = append(opts, WithExecuteStrategy(customStrategy))

	scheduler := NewScheduler(stateHandler, opts...)
	defer scheduler.Close()

	// producer
	go func() {
		ticker := time.NewTicker(20 * time.Millisecond)
		id := 1
		for {
			task := &AnalysisTask{ID: id, Payload: map[string]string{}}
			<-ticker.C
			switch id % 6 {
			case 0, 1, 2:
				task.status = StatusPending
			case 3, 4:
				task.status = StatusDownloadCompleted
				task.Payload["image"] = "image.jpg" // mock image in local memory
			case 5:
				task.status = StatusDownloadCompleted
			}
			scheduler.Submit(task)
			id++

			if id == 30 {
				fmt.Println("Producer done.")
				return
			}
		}
	}()

	select {}
}
