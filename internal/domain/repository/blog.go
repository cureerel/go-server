// internal/domain/repository/blog.go
package repository

import (
	"context"

	"github.com/cureerel/gotemplate/internal/domain/entity"
)

type BlogFilter struct {
	Page     int
	Limit    int
	Search   string
	Tags     string // partial match — ILIKE %tag%
	Status   string
	SortBy   string
	OrderDir string
	AuthorID *uint
}

type BlogRepository interface {
	Create(ctx context.Context, blog *entity.Blog) error
	GetByID(ctx context.Context, id uint) (*entity.Blog, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Blog, error)
	SlugExists(ctx context.Context, slug string) (bool, error)
	GetAll(ctx context.Context, filter BlogFilter) ([]entity.Blog, int64, error)
	GetByAuthor(ctx context.Context, authorID uint, filter BlogFilter) ([]entity.Blog, int64, error)
	Update(ctx context.Context, blog *entity.Blog) error
	Delete(ctx context.Context, id uint) error
	IncrementViews(ctx context.Context, id uint) error
	RecordView(ctx context.Context, blogID uint, visitorHash string) (bool, error)
}