package repository

import (
    "context"
    "github.com/cureerel/gotemplate/internal/domain/entity"
)

type ProductRepository interface {
    Create(ctx context.Context, product *entity.Product) error
    GetByID(ctx context.Context, id uint) (*entity.Product, error)
    GetAll(ctx context.Context, page, limit int) ([]entity.Product, int64, error)
    Update(ctx context.Context, product *entity.Product) error
    Delete(ctx context.Context, id uint) error
}