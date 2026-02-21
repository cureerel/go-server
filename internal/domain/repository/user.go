// internal/domain/repository/user_repository.go
package repository

import (
	"context"
	"github.com/cureerel/gotemplate/internal/domain/entity"
)

type UserRepository interface {
	// Legacy methods (keep for backward compatibility)
	Create(user *entity.User) error
	GetByID(id uint) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	GetAll() ([]entity.User, error)
	Update(user *entity.User) error
	Delete(id uint) error
	
	// Context methods (best practice)
	CreateContext(ctx context.Context, user *entity.User) error
	GetByIDContext(ctx context.Context, id uint) (*entity.User, error)
	GetByEmailContext(ctx context.Context, email string) (*entity.User, error)
	GetAllContext(ctx context.Context) ([]entity.User, error)
}