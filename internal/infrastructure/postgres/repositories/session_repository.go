// internal/infrastructure/postgres/repositories/session_repository.go
package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/cureerel/gotemplate/internal/domain/entity"
	"github.com/cureerel/gotemplate/internal/domain/repository"
	"github.com/cureerel/gotemplate/internal/infrastructure/postgres/models"
	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) repository.SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, session *entity.Session) error {
	m := models.SessionFromDomain(session)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	// Write back DB-generated ID
	session.ID = m.ID
	return nil
}

func (r *sessionRepository) GetByID(ctx context.Context, id string) (*entity.Session, error) {
	var m models.Session
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *sessionRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*entity.Session, error) {
	var m models.Session
	if err := r.db.WithContext(ctx).
		Where("token_hash = ? AND revoked = false", tokenHash).
		First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *sessionRepository) GetActiveByUserID(ctx context.Context, userID uint) ([]*entity.Session, error) {
	var ms []models.Session
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND revoked = false AND expires_at > ?", userID, time.Now()).
		Order("last_active DESC").
		Find(&ms).Error; err != nil {
		return nil, err
	}

	sessions := make([]*entity.Session, len(ms))
	for i, m := range ms {
		m := m // capture range var
		sessions[i] = m.ToDomain()
	}
	return sessions, nil
}

func (r *sessionRepository) Revoke(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&models.Session{}).
		Where("id = ?", id).
		Update("revoked", true).Error
}

func (r *sessionRepository) RevokeAllByUserID(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.Session{}).
		Where("user_id = ? AND revoked = false", userID).
		Update("revoked", true).Error
}

func (r *sessionRepository) UpdateLastActive(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&models.Session{}).
		Where("id = ?", id).
		Update("last_active", time.Now()).Error
}

func (r *sessionRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&models.Session{}).Error
}