package service

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
	mq PubSub
}

func (s *Service) run() error {
	const BatchSize = 100
	var now int64
	for {
		now = time.Now().Add(-time.Minute * 2).UnixMilli()
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		msgs := make([]Msg, 0, 100)
		err := s.db.WithContext(ctx).
			Where("status = ? AND utime < ?", MsgStatusInit, now).
			Offset(0).Limit(100).
			Find(&msgs).Error
		cancel()
		if err != nil {
			// TODO 判断err，如何区分处理
			continue
		}

		for _, msg := range msgs {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			err := s.sendMsg(ctx, msg)
			cancel()
			if err != nil {
				// TODO LOG
			}
			return nil
		}
	}
}

func (s *Service) sendMsg(ctx context.Context, msg Msg) error {
	err := s.mq.Pub(msg)
	if err != nil {
		return err
	}
	return s.db.WithContext(ctx).Create(&msg).Error
}

func (s *Service) Process(ctx context.Context, fn func(tx *gorm.DB) (Msg, error)) error {
	var msg Msg
	s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		msg, err = fn(tx)
		if err != nil {
			return err
		}
		now := time.Now().UnixMilli()
		msg.Utime = now
		return tx.Create(&msg).Error
	})
	err := s.mq.Pub(msg)
	if err != nil {
		return err
	}

	return nil
}

const (
	MsgStatusInit = iota
	MsgStatusSuccess
	MsgStatusFailed
)

type Msg struct {
	ID      int `gorm:"primaryKey,autoIncrement"`
	Topic   string
	Content string

	Status int8
	Utime  int64
}

type PubSub interface {
	Pub(msg Msg) error
	Sub(topic string) error
}
