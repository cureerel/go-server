package repository

import (
	"github.com/cureerel/gotemplate/internal/domain/entity"
)

type BlogRepository interface{
	Create(blog *entity.Blog) error
	GetByID(id uint) (*entity.Blog, error)
	GetBySlug(slug string) (*entity.Blog, error)
	GetAll(page, limit int) ([]entity.Blog, int64, error)
	Update(blog *entity.Blog) error
	Delete(id uint) error
	GetByAuthor(authorID uint, page, limit int) ([]entity.Blog, int64, error)
}