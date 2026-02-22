package repository

import (
    "context"
    "github.com/cureerel/gotemplate/internal/domain/entity"
)

type OrderRepository interface {
    Create(ctx context.Context, order *entity.Order) error
    GetByID(ctx context.Context, id uint) (*entity.Order, error)
    GetByUser(ctx context.Context, userID uint, page, limit int) ([]entity.Order, int64, error)
    UpdateStatus(ctx context.Context, id uint, status string) error
}