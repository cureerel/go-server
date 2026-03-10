// internal/infrastructure/postgres/repositories/coupon_usage_repository.go
package repositories

import (
	"context"
	"errors"

	"github.com/cureerel/gotemplate/internal/domain/entity"
	"github.com/cureerel/gotemplate/internal/domain/repository"
	"github.com/cureerel/gotemplate/internal/infrastructure/postgres/models"
	"gorm.io/gorm"
)

type couponUsageRepository struct{ db *gorm.DB }

func NewCouponUsageRepository(db *gorm.DB) repository.CouponUsageRepository {
	return &couponUsageRepository{db: db}
}

func (r *couponUsageRepository) Create(ctx context.Context, usage *entity.CouponUsage) error {
	m := models.CouponUsageFromDomain(usage)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	usage.ID = m.ID
	return nil
}

func (r *couponUsageRepository) GetByOrderID(ctx context.Context, orderID uint) (*entity.CouponUsage, error) {
	var m models.CouponUsage
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *couponUsageRepository) GetByCouponID(ctx context.Context, couponID uint, page, limit int) ([]entity.CouponUsage, int64, error) {
	var ms []models.CouponUsage
	var total int64
	offset := (page - 1) * limit
	q := r.db.WithContext(ctx).Model(&models.CouponUsage{}).Where("coupon_id = ?", couponID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&ms).Error; err != nil {
		return nil, 0, err
	}
	usages := make([]entity.CouponUsage, len(ms))
	for i, m := range ms {
		usages[i] = *m.ToDomain()
	}
	return usages, total, nil
}