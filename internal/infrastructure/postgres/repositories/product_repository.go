// internal/infrastructure/postgres/repositories/product_repository.go
package repositories

import (
	"context"
	"errors"

	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/domain/repository"
	"github.com/cureerel/cserver/internal/infrastructure/postgres/models"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repository.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	m := models.ProductFromDomain(product)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	product.ID = m.ID
	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id uint) (*entity.Product, error) {
	var m models.Product
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *productRepository) GetBySKU(ctx context.Context, sku string) (*entity.Product, error) {
	var m models.Product
	if err := r.db.WithContext(ctx).Where("sku = ?", sku).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *productRepository) GetAll(ctx context.Context, filter repository.ProductFilter) ([]entity.Product, int64, error) {
	var ms []models.Product
	var total int64
	offset := (filter.Page - 1) * filter.Limit

	q := r.db.WithContext(ctx).Model(&models.Product{})
	if filter.Type != "" {
		q = q.Where("type = ?", filter.Type)
	}
	if filter.IsActive != nil {
		q = q.Where("is_active = ?", *filter.IsActive)
	}
	if filter.Search != "" {
		q = q.Where("name ILIKE ?", "%"+filter.Search+"%")
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("created_at DESC").Offset(offset).Limit(filter.Limit).Find(&ms).Error; err != nil {
		return nil, 0, err
	}

	products := make([]entity.Product, len(ms))
	for i, m := range ms {
		products[i] = *m.ToDomain()
	}
	return products, total, nil
}

func (r *productRepository) Update(ctx context.Context, product *entity.Product) error {
	m := models.ProductFromDomain(product)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}

// DecrementStock reduces stock by qty. Noop for unlimited (stock=-1) products.
func (r *productRepository) DecrementStock(ctx context.Context, id uint, qty int) error {
	return r.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ? AND stock > 0", id).
		UpdateColumn("stock", gorm.Expr("stock - ?", qty)).
		Error
}
