// Package repositories
package repositories

import (
	"context"

	"sharefood/internal/entity"
)

const (
	TABLE_NAME_EXAMPLE = `example`
)

type Example interface {
	Find(ctx context.Context) ([]entity.Example, error)
	Upsert(ctx context.Context, p entity.Example) (uint64, error)
	Delete(ctx context.Context, id uint64) error
}
