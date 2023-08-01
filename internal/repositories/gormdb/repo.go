package gormdb

import (
	"GoShopper/internal/model"
	"context"
	"gorm.io/gorm"
)

type ProductGorm struct {
	db *gorm.DB
}

func New(db *gorm.DB) *ProductGorm {
	return &ProductGorm{
		db: db,
	}
}

func (r *ProductGorm) Create(ctx context.Context, s *model.Product) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *ProductGorm) Get(ctx context.Context, s *model.Product, operation string) (int64, error) {
	return 0, nil
}

func (r *ProductGorm) Delete(ctx context.Context, s *model.Product) error {

	return nil
}
