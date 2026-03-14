// internal/application/service/blog_service.go
package service

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/cureerel/gotemplate/internal/domain/entity"
	"github.com/cureerel/gotemplate/internal/domain/repository"
	"github.com/cureerel/gotemplate/pkg/apperror"
	"github.com/cureerel/gotemplate/pkg/utils"
)

type BlogService struct {
	blogRepo repository.BlogRepository
}

func NewBlogService(blogRepo repository.BlogRepository) *BlogService {
	return &BlogService{blogRepo: blogRepo}
}

// ── Input types ───────────────────────────────────────────────

type CreateBlogInput struct {
	Title         string
	Content       string
	Excerpt       string // ← added
	Status        string // ← added
	AuthorID      uint
	Tags          string
	CoverImageURL string
	CoverImageKey string
}

type UpdateBlogInput struct {
	ID            uint
	CallerID      uint
	CallerRole    string
	Title         *string
	Content       *string
	Excerpt       *string // ← added
	Status        *string
	Tags          *string
	CoverImageURL *string
	CoverImageKey *string
}

// ── CRUD ──────────────────────────────────────────────────────

func (s *BlogService) Create(ctx context.Context, in CreateBlogInput) (*entity.Blog, error) {
	slug := utils.GenerateUniqueSlug(in.Title, func(candidate string) bool {
		exists, _ := s.blogRepo.SlugExists(ctx, candidate)
		return exists
	})

	// Resolve status — default to draft if not provided or invalid
	status := entity.BlogDraft
	if in.Status == string(entity.BlogPublished) {
		status = entity.BlogPublished
	} else if in.Status == string(entity.BlogArchived) {
		status = entity.BlogArchived
	}

	blog := &entity.Blog{
		Title:         in.Title,
		Slug:          slug,
		Content:       in.Content,
		Excerpt:       in.Excerpt,
		AuthorID:      in.AuthorID,
		Tags:          in.Tags,
		Status:        status,
		CoverImageURL: in.CoverImageURL,
		CoverImageKey: in.CoverImageKey,
	}
	if err := s.blogRepo.Create(ctx, blog); err != nil {
		return nil, apperror.NewInternal(err, "failed to create blog")
	}
	return blog, nil
}

func (s *BlogService) GetByID(ctx context.Context, id uint) (*entity.Blog, error) {
	blog, err := s.blogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch blog")
	}
	if blog == nil {
		return nil, apperror.NewNotFound("blog not found")
	}
	return blog, nil
}

func (s *BlogService) GetBySlug(ctx context.Context, slug string) (*entity.Blog, error) {
	blog, err := s.blogRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch blog")
	}
	if blog == nil {
		return nil, apperror.NewNotFound("blog not found")
	}
	return blog, nil
}

func (s *BlogService) GetAll(ctx context.Context, page, limit int, search, tag string) ([]entity.Blog, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.blogRepo.GetAll(ctx, repository.BlogFilter{
		Page:   page,
		Limit:  limit,
		Status: "published",
		Search: search,
		Tags:   tag,
	})
}

func (s *BlogService) GetMine(ctx context.Context, authorID uint, page, limit int) ([]entity.Blog, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.blogRepo.GetByAuthor(ctx, authorID, repository.BlogFilter{
		Page:  page,
		Limit: limit,
	})
}

func (s *BlogService) Update(ctx context.Context, in UpdateBlogInput) (*entity.Blog, error) {
	blog, err := s.blogRepo.GetByID(ctx, in.ID)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch blog")
	}
	if blog == nil {
		return nil, apperror.NewNotFound("blog not found")
	}

	caller := &entity.User{ID: in.CallerID, Role: in.CallerRole}
	if blog.AuthorID != in.CallerID && !caller.HasRole(entity.RoleAdmin) {
		return nil, apperror.NewForbidden("you don't own this blog")
	}

	if in.Title != nil && *in.Title != "" {
		blog.Title = *in.Title
		currentSlug := blog.Slug
		blog.Slug = utils.GenerateUniqueSlug(*in.Title, func(c string) bool {
			if c == currentSlug {
				return false
			}
			exists, _ := s.blogRepo.SlugExists(ctx, c)
			return exists
		})
	}
	if in.Content != nil {
		blog.Content = *in.Content
	}
	if in.Excerpt != nil {
		blog.Excerpt = *in.Excerpt
	}
	if in.Status != nil {
		if !validBlogStatus(*in.Status) {
			return nil, apperror.NewBadRequest("invalid status — must be draft, published, or archived")
		}
		blog.Status = entity.BlogStatus(*in.Status)
	}
	if in.Tags != nil {
		blog.Tags = *in.Tags
	}
	if in.CoverImageURL != nil {
		blog.CoverImageURL = *in.CoverImageURL
	}
	if in.CoverImageKey != nil {
		blog.CoverImageKey = *in.CoverImageKey
	}

	if err := s.blogRepo.Update(ctx, blog); err != nil {
		return nil, apperror.NewInternal(err, "failed to update blog")
	}
	return blog, nil
}

func (s *BlogService) Delete(ctx context.Context, id, callerID uint, callerRole string) error {
	blog, err := s.blogRepo.GetByID(ctx, id)
	if err != nil {
		return apperror.NewInternal(err, "failed to fetch blog")
	}
	if blog == nil {
		return apperror.NewNotFound("blog not found")
	}

	caller := &entity.User{ID: callerID, Role: callerRole}
	if blog.AuthorID != callerID && !caller.HasRole(entity.RoleAdmin) {
		return apperror.NewForbidden("you don't own this blog")
	}

	return s.blogRepo.Delete(ctx, id)
}

// ── Analytics ─────────────────────────────────────────────────

func (s *BlogService) RecordView(ctx context.Context, blogID uint, ip, ua string) error {
	raw := fmt.Sprintf("%s|%s", ip, ua)
	h := sha256.Sum256([]byte(raw))
	hash := fmt.Sprintf("%x", h)
	_, err := s.blogRepo.RecordView(ctx, blogID, hash)
	return err
}

func (s *BlogService) GetStats(ctx context.Context, blogID uint) (int64, error) {
	blog, err := s.blogRepo.GetByID(ctx, blogID)
	if err != nil {
		return 0, apperror.NewInternal(err, "failed to fetch blog")
	}
	if blog == nil {
		return 0, apperror.NewNotFound("blog not found")
	}
	return blog.ViewsTotal, nil
}

// ── helpers ───────────────────────────────────────────────────

func validBlogStatus(s string) bool {
	switch s {
	case "draft", "published", "archived":
		return true
	}
	return false
}