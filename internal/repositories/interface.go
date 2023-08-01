package repository

import (
	"GoShopper/internal/model"
	"context"
)

type SocksRepository interface {
	Create(ctx context.Context, s *model.Product) error
	Get(ctx context.Context, s *model.Product, operation string) (int64, error)
	Delete(ctx context.Context, s *model.Product) error
}
