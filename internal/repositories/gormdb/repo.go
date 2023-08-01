package gormdb

import (
	"GoShopper/internal/model"
	"context"
	"gorm.io/gorm"
)

type SocksGorm struct {
	db *gorm.DB
}

func New(db *gorm.DB) *SocksGorm {
	return &SocksGorm{
		db: db,
	}
}

func (r *SocksGorm) Create(ctx context.Context, s *model.Product) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *SocksGorm) Get(ctx context.Context, s *model.Product, operation string) (int64, error) {
	return 0, nil
}

func (r *SocksGorm) Delete(ctx context.Context, s *model.Product) error {

	return nil
}
