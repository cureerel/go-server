package persistence

import (
	"github.com/cureerel/gotemplate/internal/domain/entity"
	"github.com/cureerel/gotemplate/internal/domain/repository"

	"gorm.io/gorm"
)

type BlogRepositoryImpl struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) repository.BlogRepository {
	return &BlogRepositoryImpl{db: db}
}

func (r *BlogRepositoryImpl) Create(blog *entity.Blog) error {
	return r.db.Create(blog).Error
}

func (r *BlogRepositoryImpl) GetByID(id uint) (*entity.Blog, error) {
	var blog entity.Blog
	if err := r.db.First(&blog, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &blog, nil
}

func (r *BlogRepositoryImpl) GetBySlug(slug string) (*entity.Blog, error) {
	var blog entity.Blog
	if err := r.db.Where("slug = ?", slug).First(&blog).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &blog, nil
}

func (r *BlogRepositoryImpl) GetAll(page, limit int) ([]entity.Blog, int64, error) {
	var blogs []entity.Blog
	var total int64
	
	offset := (page - 1) * limit
	
	if err := r.db.Model(&entity.Blog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	if err := r.db.Order("created_at desc").Limit(limit).Offset(offset).Find(&blogs).Error; err != nil {
		return nil, 0, err
	}
	
	return blogs, total, nil
}

func (r *BlogRepositoryImpl) Update(blog *entity.Blog) error {
	return r.db.Save(blog).Error
}

func (r *BlogRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entity.Blog{}, id).Error
}

func (r *BlogRepositoryImpl) GetByAuthor(authorID uint, page, limit int) ([]entity.Blog, int64, error) {
	var blogs []entity.Blog
	var total int64
	
	offset := (page - 1) * limit
	
	if err := r.db.Model(&entity.Blog{}).Where("author_id = ?", authorID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	if err := r.db.Where("author_id = ?", authorID).Order("created_at desc").Limit(limit).Offset(offset).Find(&blogs).Error; err != nil {
		return nil, 0, err
	}
	
	return blogs, total, nil
}