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

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) repository.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) SaveRefreshToken(ctx context.Context, token *entity.RefreshToken) error {
	m := models.RefreshTokenFromDomain(token)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	token.ID = m.ID
	return nil
}

func (r *authRepository) GetRefreshToken(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	var m models.RefreshToken
	if err := r.db.WithContext(ctx).Where("token_hash = ? AND revoked = false", tokenHash).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *authRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	return r.db.WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("token_hash = ?", tokenHash).
		Update("revoked", true).Error
}

func (r *authRepository) BlacklistToken(ctx context.Context, tokenHash string, expiresAt time.Time) error {
	m := &models.BlacklistedToken{
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *authRepository) IsTokenBlacklisted(ctx context.Context, tokenHash string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.BlacklistedToken{}).
		Where("token_hash = ? AND expires_at > ?", tokenHash, time.Now()).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *authRepository) RevokeAllRefreshTokens(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked = false", userID).
		Update("revoked", true).Error
}
