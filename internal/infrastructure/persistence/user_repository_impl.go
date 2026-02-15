package persistence

import (
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

func (r *UserRepositoryImpl) Create(user *entity.User) error {
    return r.db.Create(user).Error
}

func (r *UserRepositoryImpl) FindByID(id int) (*entity.User, error) {
    var user entity.User
    err := r.db.First(&user, id).Error
    return &user, err
}

func (r *UserRepositoryImpl) FindAll() ([]*entity.User, error) {
    var users []*entity.User
    err := r.db.Find(&users).Error
    return users, err
}