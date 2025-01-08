package service

import (
	"context"

	"gorm.io/gorm"
)

type Executor interface {
	Execute(ctx context.Context, db *gorm.DB, table string) error
}
