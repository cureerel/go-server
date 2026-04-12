// internal/infrastructure/postgres/repositories/coupon_repository.go
package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/domain/repository"
	"github.com/cureerel/cserver/internal/infrastructure/postgres/models"
	"gorm.io/gorm"
)

type couponRepository struct{ db *gorm.DB }

func NewCouponRepository(db *gorm.DB) repository.CouponRepository {
	return &couponRepository{db: db}
}

func (r *couponRepository) Create(ctx context.Context, c *entity.Coupon) error {
	m := models.CouponFromDomain(c)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	c.ID = m.ID
	return nil
}

func (r *couponRepository) GetByID(ctx context.Context, id uint) (*entity.Coupon, error) {
	var m models.Coupon
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *couponRepository) GetByCode(ctx context.Context, code string) (*entity.Coupon, error) {
	var m models.Coupon
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *couponRepository) GetAll(ctx context.Context, page, limit int, status string) ([]entity.Coupon, int64, error) {
	var ms []models.Coupon
	var total int64
	offset := (page - 1) * limit
	q := r.db.WithContext(ctx).Model(&models.Coupon{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	q = q.Order("created_at DESC")
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(offset).Limit(limit).Find(&ms).Error; err != nil {
		return nil, 0, err
	}
	coupons := make([]entity.Coupon, len(ms))
	for i, m := range ms {
		coupons[i] = *m.ToDomain()
	}
	return coupons, total, nil
}

func (r *couponRepository) UpdateStatus(ctx context.Context, id uint, status string, approvedBy uint) error {
	now := time.Now()
	updates := map[string]any{
		"status":      status,
		"approved_at": now,
	}
	if approvedBy != 0 {
		updates["approved_by"] = approvedBy
	}
	return r.db.WithContext(ctx).Model(&models.Coupon{}).Where("id = ?", id).Updates(updates).Error
}

func (r *couponRepository) IncrementUsed(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&models.Coupon{}).
		Where("id = ?", id).UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error
}