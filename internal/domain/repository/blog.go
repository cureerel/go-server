package repository

import (
    "context"
    "github.com/cureerel/gotemplate/internal/domain/entity"
)

type BlogFilter struct {
    Page     int
    Limit    int
    SortBy   string
    OrderDir string
    Search   string
}

type BlogRepository interface {
    Create(ctx context.Context, blog *entity.Blog) error
    GetByID(ctx context.Context, id uint) (*entity.Blog, error)
    GetBySlug(ctx context.Context, slug string) (*entity.Blog, error)
    GetAll(ctx context.Context, filter BlogFilter) ([]entity.Blog, int64, error)
    Update(ctx context.Context, blog *entity.Blog) error
    Delete(ctx context.Context, id uint) error
    GetByAuthor(ctx context.Context, authorID uint, filter BlogFilter) ([]entity.Blog, int64, error)
}