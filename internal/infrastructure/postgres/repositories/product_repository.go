package repositories

import (
    "context"
    "errors"

    "github.com/cureerel/gotemplate/internal/domain/entity"
    "github.com/cureerel/gotemplate/internal/domain/repository"
    "github.com/cureerel/gotemplate/internal/infrastructure/postgres/models"
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

func (r *productRepository) GetAll(ctx context.Context, page, limit int) ([]entity.Product, int64, error) {
    var ms []models.Product
    var total int64
    offset := (page - 1) * limit

    if err := r.db.WithContext(ctx).Model(&models.Product{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }
    if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&ms).Error; err != nil {
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