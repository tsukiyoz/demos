package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

var ErrProcessFailed = fmt.Errorf("process failed")

const (
	TypeEmailDelivery = "email:deliver"
	TypeImageResize   = "image:resize"
)

type EmailDeliveryPayload struct {
	UserID     int
	TemplateID string
	SendAt     time.Time
}

type ImageResizePayload struct {
	SourceURL string
}

func NewEmailDeliveryTask(userID int, tmplID string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailDeliveryPayload{
		UserID:     userID,
		TemplateID: tmplID,
		SendAt:     time.Now(),
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeEmailDelivery, payload), nil
}

func NewImageResizeTask(src string) (*asynq.Task, error) {
	payload, err := json.Marshal(ImageResizePayload{
		SourceURL: src,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeImageResize, payload, asynq.MaxRetry(10), asynq.Timeout(20*time.Minute)), nil
}

func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
	fmt.Printf("got task [%s] at %v\n", t.Type(), time.Now().String())
	var payload EmailDeliveryPayload
	err := json.Unmarshal(t.Payload(), &payload)
	if err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	fmt.Printf("handled email task %v\n", payload)

	return nil
}

type ImageProcessor struct{}

func (i *ImageProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	fmt.Printf("got task [%s] at %v\n", t.Type(), time.Now().String())
	var payload ImageResizePayload
	err := json.Unmarshal(t.Payload(), &payload)
	if err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Resizing image: src=%s", payload.SourceURL)
	// mock failed
	if time.Now().Second()%3 != 0 {
		fmt.Printf("failed to process task: %v\n", errors.New("mock failure"))
		return ErrProcessFailed
	}
	fmt.Printf("handled image task %v\n", payload)
	return nil
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}
