// internal/infrastructure/postgres/repositories/payment_repository.go
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

type paymentRepository struct{ db *gorm.DB }

func NewPaymentRepository(db *gorm.DB) repository.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, p *entity.Payment) error {
	m := models.PaymentFromDomain(p)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *paymentRepository) GetByID(ctx context.Context, id string) (*entity.Payment, error) {
	var m models.Payment
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *paymentRepository) GetByOrderID(ctx context.Context, orderID uint) (*entity.Payment, error) {
	var m models.Payment
	if err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *paymentRepository) GetAll(ctx context.Context, filter repository.PaymentFilter) ([]entity.Payment, int64, error) {
	var ms []models.Payment
	var total int64
	offset := (filter.Page - 1) * filter.Limit

	q := r.db.WithContext(ctx).Model(&models.Payment{})
	if filter.UserID != nil {
		q = q.Where("user_id = ?", *filter.UserID)
	}
	if filter.OrderID != nil {
		q = q.Where("order_id = ?", *filter.OrderID)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	q = q.Order("created_at DESC")

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(offset).Limit(filter.Limit).Find(&ms).Error; err != nil {
		return nil, 0, err
	}
	payments := make([]entity.Payment, len(ms))
	for i, m := range ms {
		payments[i] = *m.ToDomain()
	}
	return payments, total, nil
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, id string, status entity.PaymentStatus) error {
	return r.db.WithContext(ctx).Model(&models.Payment{}).
		Where("id = ?", id).Update("status", string(status)).Error
}

func (r *paymentRepository) MarkRefunded(ctx context.Context, id string, refundID string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.Payment{}).
		Where("id = ?", id).Updates(map[string]any{
		"status":      string(entity.PaymentRefunded),
		"refund_id":   refundID,
		"refunded_at": now,
	}).Error
}

func (r *paymentRepository) SetProviderTxnID(ctx context.Context, id string, txnID string) error {
	return r.db.WithContext(ctx).Model(&models.Payment{}).
		Where("id = ?", id).Update("provider_txn_id", txnID).Error
}
