package persistence

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"


	"github.com/cureerel/gotemplate/internal/domain/repository"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    string    `gorm:"index"`
	TokenHash string    `gorm:"uniqueIndex"`
	ExpiresAt time.Time
	Revoked   bool `gorm:"default:false"`
	CreatedAt time.Time
}

func NewAuthRepository(db *gorm.DB) repository.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) StoreRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time) error {
	return r.db.WithContext(ctx).Create(&RefreshToken{
		UserID:    userID,
		TokenHash: hashToken(token),
		ExpiresAt: expiresAt,
	}).Error
}

func (r *authRepository) ValidateRefreshToken(ctx context.Context, token string) (string, error) {
	var rt RefreshToken
	err := r.db.WithContext(ctx).
		Where("token_hash = ? AND revoked = ? AND expires_at > ?", 
			hashToken(token), false, time.Now()).
		First(&rt).Error
	if err != nil {
		return "", err
	}
	return rt.UserID, nil
}

func (r *authRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Model(&RefreshToken{}).
		Where("token_hash = ?", hashToken(token)).
		Update("revoked", true).Error
}

// hashToken returns SHA256 hash of token
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}