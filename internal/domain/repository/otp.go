// internal/domain/repository/otp.go
package repository

import (
	"context"
	"github.com/cureerel/gotemplate/internal/domain/entity"
)

type OTPRepository interface {
	// Create inserts a new OTP record.
	Create(ctx context.Context, otp *entity.OTP) error

	// GetLatestByEmail returns the most recent unused, unexpired OTP
	// for the given email and type.
	GetLatestByEmail(ctx context.Context, email, otpType string) (*entity.OTP, error)

	// MarkUsed marks an OTP as used so it cannot be reused.
	MarkUsed(ctx context.Context, id uint) error

	// DeleteExpiredByEmail cleans up old OTPs for an email.
	// Call before creating a new OTP to avoid DB bloat.
	DeleteExpiredByEmail(ctx context.Context, email, otpType string) error
}