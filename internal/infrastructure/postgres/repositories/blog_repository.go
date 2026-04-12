// internal/infrastructure/postgres/repositories/blog_repository.go
package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/domain/repository"
	"github.com/cureerel/cserver/internal/infrastructure/postgres/models"
	"gorm.io/gorm"
)

type blogRepository struct{ db *gorm.DB }

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

func (r *blogRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Blog{}).Where("slug = ?", slug).Count(&count).Error
	return count > 0, err
}

func (r *blogRepository) GetAll(ctx context.Context, filter repository.BlogFilter) ([]entity.Blog, int64, error) {
	var ms []models.Blog
	var total int64
	offset := (filter.Page - 1) * filter.Limit
	q := r.db.WithContext(ctx).Model(&models.Blog{})

	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	if filter.Search != "" {
		q = q.Where("title ILIKE ?", "%"+filter.Search+"%")
	}
	// Tags stored as comma-separated string — partial ILIKE match covers any tag
	if filter.Tags != "" {
		q = q.Where("tags ILIKE ?", "%"+filter.Tags+"%")
	}
	if filter.AuthorID != nil {
		q = q.Where("author_id = ?", *filter.AuthorID)
	}

	if filter.SortBy != "" {
		dir := "ASC"
		if filter.OrderDir == "desc" || filter.OrderDir == "DESC" {
			dir = "DESC"
		}
		q = q.Order(filter.SortBy + " " + dir)
	} else {
		q = q.Order("created_at DESC")
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

func (r *blogRepository) GetByAuthor(ctx context.Context, authorID uint, filter repository.BlogFilter) ([]entity.Blog, int64, error) {
	filter.AuthorID = &authorID
	return r.GetAll(ctx, filter)
}

func (r *blogRepository) Update(ctx context.Context, blog *entity.Blog) error {
	return r.db.WithContext(ctx).Save(models.BlogFromDomain(blog)).Error
}

func (r *blogRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Blog{}, id).Error
}

func (r *blogRepository) IncrementViews(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&models.Blog{}).
		Where("id = ?", id).
		UpdateColumn("views_total", gorm.Expr("views_total + 1")).Error
}

// RecordView inserts a unique (blog_id, visitor_hash, viewed_date) row.
// Returns true only on first view of the day for that visitor.
// Increments views_total on the blog row when a new view is recorded.
func (r *blogRepository) RecordView(ctx context.Context, blogID uint, visitorHash string) (bool, error) {
	today := time.Now().UTC().Format("2006-01-02")
	result := r.db.WithContext(ctx).Exec(`
		INSERT INTO blog_views (blog_id, visitor_hash, viewed_date, created_at)
		VALUES (?, ?, ?, NOW())
		ON CONFLICT (blog_id, visitor_hash, viewed_date) DO NOTHING`,
		blogID, visitorHash, today,
	)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected > 0 {
		_ = r.IncrementViews(ctx, blogID)
		return true, nil
	}
	return false, nil
}

func (r *blogRepository) AddCoAuthor(ctx context.Context, blogID, userID uint) error {
	if blogID == 0 || userID == 0 {
		return nil
	}
	var b models.Blog
	if err := r.db.WithContext(ctx).Select("author_id").First(&b, blogID).Error; err != nil {
		return err
	}
	if b.AuthorID == userID {
		return nil
	}
	return r.db.WithContext(ctx).Exec(`
		INSERT INTO blog_authors (blog_id, user_id, role, created_at)
		VALUES (?, ?, 'co_author', NOW())
		ON CONFLICT (blog_id, user_id) DO NOTHING`, blogID, userID).Error
}

func (r *blogRepository) ListCoAuthors(ctx context.Context, blogID uint) ([]entity.BlogAuthor, error) {
	var rows []models.BlogAuthor
	if err := r.db.WithContext(ctx).Where("blog_id = ?", blogID).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]entity.BlogAuthor, len(rows))
	for i, row := range rows {
		out[i] = entity.BlogAuthor{BlogID: row.BlogID, UserID: row.UserID, Role: row.Role}
	}
	return out, nil
}

func (r *blogRepository) UserCanEditBlog(ctx context.Context, blogID, userID uint) (bool, error) {
	var b models.Blog
	if err := r.db.WithContext(ctx).Select("author_id").First(&b, blogID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	if b.AuthorID == userID {
		return true, nil
	}
	var n int64
	err := r.db.WithContext(ctx).Model(&models.BlogAuthor{}).
		Where("blog_id = ? AND user_id = ?", blogID, userID).Count(&n).Error
	return n > 0, err
}

func (r *blogRepository) ListByStatus(ctx context.Context, status string, page, limit int) ([]entity.Blog, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	var total int64
	var ms []models.Blog
	q := r.db.WithContext(ctx).Model(&models.Blog{}).Where("status = ?", status)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	if err := q.Order("submitted_for_review_at DESC NULLS LAST").Offset(offset).Limit(limit).Find(&ms).Error; err != nil {
		return nil, 0, err
	}
	out := make([]entity.Blog, len(ms))
	for i := range ms {
		out[i] = *ms[i].ToDomain()
	}
	return out, total, nil
}