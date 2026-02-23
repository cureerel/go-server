package service

import (
    "context"
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

func (s *BlogService) Create(ctx context.Context, input CreateBlogInput) (*entity.Blog, error) {
    slug := utils.GenerateSlug(input.Title)
    blog := &entity.Blog{
        Title:    input.Title,
        Slug:     slug,
        Content:  input.Content,
        AuthorID: input.AuthorID,
        Tags:     input.Tags,
        Status:   "draft",
    }
    if err := s.blogRepo.Create(ctx, blog); err != nil {
        return nil, err
    }
    return blog, nil
}

func (s *BlogService) GetByID(ctx context.Context, id uint) (*entity.Blog, error) {
    return s.blogRepo.GetByID(ctx, id)
}

func (s *BlogService) GetBySlug(ctx context.Context, slug string) (*entity.Blog, error) {
    return s.blogRepo.GetBySlug(ctx, slug)
}

func (s *BlogService) GetAll(ctx context.Context, page, limit int) ([]entity.Blog, int64, error) {
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 10
    }
    filter := repository.BlogFilter{
        Page:  page,
        Limit: limit,
    }
    return s.blogRepo.GetAll(ctx, filter)
}

func (s *BlogService) Update(ctx context.Context, input UpdateBlogInput) (*entity.Blog, error) {
    blog, err := s.blogRepo.GetByID(ctx, input.ID)
    if err != nil {
        return nil, err
    }
    if blog == nil {
        return nil, errors.New("blog not found")
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

    if err := s.blogRepo.Update(ctx, blog); err != nil {
        return nil, err
    }
    return blog, nil
}

func (s *BlogService) Delete(ctx context.Context, id uint) error {
    return s.blogRepo.Delete(ctx, id)
}

func isValidStatus(status string) bool {
    validStatuses := map[string]bool{
        "draft":     true,
        "published": true,
        "archived":  true,
    }
    return validStatuses[status]
}