// internal/domain/repository/user.go
package repository

import (
	"context"
	"github.com/cureerel/cserver/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uint) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetAll(ctx context.Context, page, limit int) ([]entity.User, int64, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uint) error

	// UpdateRole
	UpdateRole(ctx context.Context, id uint, role string) error

	// UpdateVerified after OTP confirmation.
	UpdateVerified(ctx context.Context, id uint, verified bool) error

	// UpdatePassword
	UpdatePassword(ctx context.Context, id uint, passwordHash string) error
}
