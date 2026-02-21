package persistence

import (
	"context"
	"github.com/cureerel/gotemplate/internal/domain/entity"
	"github.com/cureerel/gotemplate/internal/domain/repository"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepositoryImpl{db: db}
}

// Legacy methods (no context)
func (r *UserRepositoryImpl) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepositoryImpl) GetByID(id uint) (*entity.User, error) {
	return r.GetByIDContext(context.Background(), id)
}

func (r *UserRepositoryImpl) GetByEmail(email string) (*entity.User, error) {
	return r.GetByEmailContext(context.Background(), email)
}

func (r *UserRepositoryImpl) GetAll() ([]entity.User, error) {
	return r.GetAllContext(context.Background())
}

func (r *UserRepositoryImpl) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entity.User{}, id).Error
}

// Context methods (required by interface)
func (r *UserRepositoryImpl) CreateContext(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepositoryImpl) GetByIDContext(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetByEmailContext(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetAllContext(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}