// internal/infrastructure/postgres/repositories/otp_repository.go
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

type otpRepository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) repository.OTPRepository {
	return &otpRepository{db: db}
}

func (r *otpRepository) Create(ctx context.Context, otp *entity.OTP) error {
	m := models.OTPFromDomain(otp)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	otp.ID = m.ID
	return nil
}

func (r *otpRepository) GetLatestByEmail(ctx context.Context, email, otpType string) (*entity.OTP, error) {
	var m models.OTP
	err := r.db.WithContext(ctx).
		Where("email = ? AND type = ? AND used = false AND expires_at > ?", email, otpType, time.Now()).
		Order("created_at DESC").
		First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *otpRepository) MarkUsed(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&models.OTP{}).
		Where("id = ?", id).
		Update("used", true).Error
}

func (r *otpRepository) DeleteExpiredByEmail(ctx context.Context, email, otpType string) error {
	return r.db.WithContext(ctx).
		Where("email = ? AND type = ? AND (expires_at < ? OR used = true)", email, otpType, time.Now()).
		Delete(&models.OTP{}).Error
}