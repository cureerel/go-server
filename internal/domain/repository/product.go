// internal/domain/repository/product.go
package repository

import (
	"context"

	"github.com/cureerel/cserver/internal/domain/entity"
)

type ProductFilter struct {
	Page     int
	Limit    int
	Type     string
	IsActive *bool
	Search   string
}

type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	GetByID(ctx context.Context, id uint) (*entity.Product, error)
	GetBySKU(ctx context.Context, sku string) (*entity.Product, error)
	GetAll(ctx context.Context, filter ProductFilter) ([]entity.Product, int64, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id uint) error
	DecrementStock(ctx context.Context, id uint, qty int) error
}
