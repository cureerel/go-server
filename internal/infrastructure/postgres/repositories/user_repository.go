// internal/infrastructure/postgres/repositories/user_repository.go
package repositories

import (
	"context"
	"errors"

	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/domain/repository"
	"github.com/cureerel/cserver/internal/infrastructure/postgres/models"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	m := models.UserFromDomain(user)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	user.ID = m.ID
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	var m models.User
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var m models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *userRepository) GetAll(ctx context.Context, page, limit int) ([]entity.User, int64, error) {
	var ms []models.User
	var total int64
	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&ms).Error; err != nil {
		return nil, 0, err
	}

	users := make([]entity.User, len(ms))
	for i, m := range ms {
		users[i] = *m.ToDomain()
	}
	return users, total, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	m := models.UserFromDomain(user)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (r *userRepository) UpdateRole(ctx context.Context, id uint, role string) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("role", role).Error
}

func (r *userRepository) UpdateVerified(ctx context.Context, id uint, verified bool) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("is_verified", verified).Error
}

func (r *userRepository) UpdatePassword(ctx context.Context, id uint, passwordHash string) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("password_hash", passwordHash).Error
}
