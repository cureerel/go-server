// internal/infrastructure/postgres/repositories/upgrade_request_repository.go
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

type upgradeRequestRepository struct{ db *gorm.DB }

func NewUpgradeRequestRepository(db *gorm.DB) repository.UpgradeRequestRepository {
	return &upgradeRequestRepository{db: db}
}

func (r *upgradeRequestRepository) Create(ctx context.Context, req *entity.UpgradeRequest) error {
	m := models.UpgradeRequestFromDomain(req)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	req.ID = m.ID
	return nil
}

func (r *upgradeRequestRepository) GetByID(ctx context.Context, id uint) (*entity.UpgradeRequest, error) {
	var m models.UpgradeRequest
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *upgradeRequestRepository) GetPending(ctx context.Context, page, limit int) ([]entity.UpgradeRequest, int64, error) {
	var ms []models.UpgradeRequest
	var total int64
	offset := (page - 1) * limit
	q := r.db.WithContext(ctx).Model(&models.UpgradeRequest{}).
		Where("status = ?", entity.UpgradeStatusPending).
		Order("created_at ASC")
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(offset).Limit(limit).Find(&ms).Error; err != nil {
		return nil, 0, err
	}
	reqs := make([]entity.UpgradeRequest, len(ms))
	for i, m := range ms {
		reqs[i] = *m.ToDomain()
	}
	return reqs, total, nil
}

func (r *upgradeRequestRepository) GetByUser(ctx context.Context, userID uint) (*entity.UpgradeRequest, error) {
	var m models.UpgradeRequest
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, entity.UpgradeStatusPending).
		First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *upgradeRequestRepository) Review(ctx context.Context, id uint, status string, reviewedBy uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.UpgradeRequest{}).Where("id = ?", id).Updates(map[string]any{
		"status":      status,
		"reviewed_by": reviewedBy,
		"reviewed_at": now,
	}).Error
}