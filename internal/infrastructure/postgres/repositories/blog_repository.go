package repositories

import (
    "context"
    "errors"

    "github.com/cureerel/gotemplate/internal/domain/entity"
    "github.com/cureerel/gotemplate/internal/domain/repository"
    "github.com/cureerel/gotemplate/internal/infrastructure/postgres/models"
    "gorm.io/gorm"
)

type blogRepository struct {
    db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) repository.BlogRepository {
    return &blogRepository{db: db}
}

func (r *blogRepository) Create(ctx context.Context, blog *entity.Blog) error {
    m := models.BlogFromDomain(blog)
    if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
        return err
    }
    blog.ID = m.ID
    return nil
}

func (r *blogRepository) GetByID(ctx context.Context, id uint) (*entity.Blog, error) {
    var m models.Blog
    if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return m.ToDomain(), nil
}

func (r *blogRepository) GetBySlug(ctx context.Context, slug string) (*entity.Blog, error) {
    var m models.Blog
    if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&m).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return m.ToDomain(), nil
}

func (r *blogRepository) GetAll(ctx context.Context, filter repository.BlogFilter) ([]entity.Blog, int64, error) {
    var ms []models.Blog
    var total int64

    offset := (filter.Page - 1) * filter.Limit
    q := r.db.WithContext(ctx).Model(&models.Blog{})

    if filter.Search != "" {
        q = q.Where("title ILIKE ?", "%"+filter.Search+"%")
    }
    if filter.SortBy != "" {
        order := filter.SortBy
        if filter.OrderDir != "" {
            order += " " + filter.OrderDir
        }
        q = q.Order(order)
    }

    if err := q.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    if err := q.Offset(offset).Limit(filter.Limit).Find(&ms).Error; err != nil {
        return nil, 0, err
    }

    blogs := make([]entity.Blog, len(ms))
    for i, m := range ms {
        blogs[i] = *m.ToDomain()
    }
    return blogs, total, nil
}

func (r *blogRepository) Update(ctx context.Context, blog *entity.Blog) error {
    m := models.BlogFromDomain(blog)
    return r.db.WithContext(ctx).Save(m).Error
}

func (r *blogRepository) Delete(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&models.Blog{}, id).Error
}

func (r *blogRepository) GetByAuthor(ctx context.Context, authorID uint, filter repository.BlogFilter) ([]entity.Blog, int64, error) {
    var ms []models.Blog
    var total int64

    offset := (filter.Page - 1) * filter.Limit
    q := r.db.WithContext(ctx).Model(&models.Blog{}).Where("author_id = ?", authorID)

    if err := q.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    if err := q.Offset(offset).Limit(filter.Limit).Find(&ms).Error; err != nil {
        return nil, 0, err
    }

    blogs := make([]entity.Blog, len(ms))
    for i, m := range ms {
        blogs[i] = *m.ToDomain()
    }
    return blogs, total, nil
}