
package repository

import "github.com/cureerel/gotemplate/internal/domain/entity"

type UserRepository interface {
    Create(user *entity.User) error
    FindByID(id int) (*entity.User, error)
    FindAll() ([]*entity.User, error)
}