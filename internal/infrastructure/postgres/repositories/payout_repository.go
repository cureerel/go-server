// internal/infrastructure/postgres/repositories/payout_repository.go
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

type payoutRepository struct{ db *gorm.DB }

func NewPayoutRepository(db *gorm.DB) repository.PayoutRepository {
	return &payoutRepository{db: db}
}

func (r *payoutRepository) Create(ctx context.Context, p *entity.Payout) error {
	m := models.PayoutFromDomain(p)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	p.ID = m.ID
	return nil
}

func (r *payoutRepository) GetByID(ctx context.Context, id uint) (*entity.Payout, error) {
	var m models.Payout
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *payoutRepository) GetAll(ctx context.Context, page, limit int, status string) ([]entity.Payout, int64, error) {
	var ms []models.Payout
	var total int64
	offset := (page - 1) * limit
	q := r.db.WithContext(ctx).Model(&models.Payout{})
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
	payouts := make([]entity.Payout, len(ms))
	for i, m := range ms {
		payouts[i] = *m.ToDomain()
	}
	return payouts, total, nil
}

func (r *payoutRepository) GetByRecipient(ctx context.Context, recipientID uint, page, limit int) ([]entity.Payout, int64, error) {
	var ms []models.Payout
	var total int64
	offset := (page - 1) * limit
	q := r.db.WithContext(ctx).Model(&models.Payout{}).
		Where("recipient_id = ?", recipientID).
		Order("created_at DESC")
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(offset).Limit(limit).Find(&ms).Error; err != nil {
		return nil, 0, err
	}
	payouts := make([]entity.Payout, len(ms))
	for i, m := range ms {
		payouts[i] = *m.ToDomain()
	}
	return payouts, total, nil
}

func (r *payoutRepository) MarkPaid(ctx context.Context, id uint, paidBy uint, reference string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.Payout{}).Where("id = ?", id).Updates(map[string]any{
		"status":    entity.PayoutStatusPaid,
		"paid_by":   paidBy,
		"paid_at":   now,
		"reference": reference,
	}).Error
}