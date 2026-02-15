package service

import (
	"errors"
	"github.com/cureerel/gotemplate/internal/domain/entity"
	"github.com/cureerel/gotemplate/internal/domain/repository"
	"github.com/cureerel/gotemplate/pkg/utils"

)

type BlogService struct {
	blogRepo repository.BlogRepository
}

func NewBlogService(blogRepo repository.BlogRepository) *BlogService {
	return &BlogService{blogRepo: blogRepo}
}

type CreateBlogInput struct {
	Title    string
	Content  string
	AuthorID uint
	Tags     string
}

type UpdateBlogInput struct {
	ID      uint
	Title   *string
	Content *string
	Status  *string
	Tags    *string
}

func (s *BlogService) Create(input CreateBlogInput) (*entity.Blog, error) {
	slug := utils.GenerateSlug(input.Title)
	
	blog := &entity.Blog{
		Title:    input.Title,
		Slug:     slug,
		Content:  input.Content,
		AuthorID: input.AuthorID,
		Tags:     input.Tags,
		Status:   "draft",
	}
	
	if err := s.blogRepo.Create(blog); err != nil {
		return nil, err
	}
	return blog, nil
}

func (s *BlogService) GetByID(id uint) (*entity.Blog, error) {
	return s.blogRepo.GetByID(id)
}

func (s *BlogService) GetBySlug(slug string) (*entity.Blog, error) {
	return s.blogRepo.GetBySlug(slug)
}

func (s *BlogService) GetAll(page, limit int) ([]entity.Blog, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.blogRepo.GetAll(page, limit)
}

func (s *BlogService) Update(input UpdateBlogInput) (*entity.Blog, error) {
	blog, err := s.blogRepo.GetByID(input.ID)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		blog.Title = *input.Title
		blog.Slug = utils.GenerateSlug(*input.Title)
	}
	if input.Content != nil {
		blog.Content = *input.Content
	}
	if input.Status != nil {
		if !isValidStatus(*input.Status) {
			return nil, errors.New("invalid status")
		}
		blog.Status = *input.Status
	}
	if input.Tags != nil {
		blog.Tags = *input.Tags
	}

	if err := s.blogRepo.Update(blog); err != nil {
		return nil, err
	}
	return blog, nil
}

func (s *BlogService) Delete(id uint) error {
	return s.blogRepo.Delete(id)
}

func isValidStatus(status string) bool {
	validStatuses := map[string]bool{
		"draft":     true,
		"published": true,
		"archived":  true,
	}
	return validStatuses[status]
}


