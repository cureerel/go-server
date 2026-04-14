// internal/application/service/blog_service.go
package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/domain/repository"
	"github.com/cureerel/cserver/pkg/apperror"
	"github.com/cureerel/cserver/pkg/utils"
)

type BlogService struct {
	blogRepo       repository.BlogRepository
	coinRepo       repository.CoinRepository
	membershipRepo repository.MembershipRepository
}

func NewBlogService(
	blogRepo repository.BlogRepository,
	coinRepo repository.CoinRepository,
	membershipRepo repository.MembershipRepository,
) *BlogService {
	return &BlogService{
		blogRepo:       blogRepo,
		coinRepo:       coinRepo,
		membershipRepo: membershipRepo,
	}
}

// Input types

type CreateBlogInput struct {
	Title        string
	Content      string
	Keyword      string
	Tag          string
	Excerpt      string
	Thumbnail    string
	ThumbnailKey string
	Status       string
	AccessType   string
	CoinPrice    int64
}

type UpdateBlogInput struct {
	ID           uint
	CallerID     uint
	CallerRole   string
	Title        *string
	Content      *string
	Keyword      *string
	Tag          *string
	Excerpt      *string
	Thumbnail    *string
	ThumbnailKey *string
	Status       *string
	AccessType   *string
	CoinPrice    *int64
}

// CRUD

func (s *BlogService) Create(ctx context.Context, in CreateBlogInput) (*entity.Blog, error) {
	slug := utils.GenerateUniqueSlug(in.Title, func(candidate string) bool {
		exists, _ := s.blogRepo.SlugExists(ctx, candidate)
		return exists
	})

	// Admin publishes directly; default draft
	status := entity.BlogDraft
	if in.Status == string(entity.BlogPublished) {
		status = entity.BlogPublished
	} else if in.Status == string(entity.BlogArchived) {
		status = entity.BlogArchived
	}

	at := entity.AccessFree
	switch strings.TrimSpace(in.AccessType) {
	case "member":
		at = entity.AccessMember
	case "paid_coins":
		at = entity.AccessPaidCoins
	}
	coinPrice := in.CoinPrice
	if at != entity.AccessPaidCoins {
		coinPrice = 0
	}

	blog := &entity.Blog{
		Title:        in.Title,
		Slug:         slug,
		Content:      in.Content,
		Keyword:      in.Keyword,
		Tag:          in.Tag,
		Excerpt:      in.Excerpt,
		Thumbnail:    in.Thumbnail,
		ThumbnailKey: in.ThumbnailKey,
		Status:       status,
		AccessType:   at,
		CoinPrice:    coinPrice,
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
		Tag:    tag,
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

	// Admin only
	caller := &entity.User{ID: in.CallerID, Role: in.CallerRole}
	if !caller.HasRole(entity.RoleAdmin) {
		return nil, apperror.NewForbidden("only admin can edit blogs")
	}

	if in.Title != nil && *in.Title != "" {
		blog.Title = *in.Title
		cur := blog.Slug
		blog.Slug = utils.GenerateUniqueSlug(*in.Title, func(c string) bool {
			if c == cur {
				return false
			}
			exists, _ := s.blogRepo.SlugExists(ctx, c)
			return exists
		})
	}
	if in.Content != nil {
		blog.Content = *in.Content
	}
	if in.Keyword != nil {
		blog.Keyword = *in.Keyword
	}
	if in.Tag != nil {
		blog.Tag = *in.Tag
	}
	if in.Excerpt != nil {
		blog.Excerpt = *in.Excerpt
	}
	if in.Thumbnail != nil {
		blog.Thumbnail = *in.Thumbnail
	}
	if in.ThumbnailKey != nil {
		blog.ThumbnailKey = *in.ThumbnailKey
	}
	if in.Status != nil {
		if !validBlogStatus(*in.Status) {
			return nil, apperror.NewBadRequest("invalid blog status")
		}
		blog.Status = entity.BlogStatus(*in.Status)
	}
	if in.AccessType != nil {
		switch strings.TrimSpace(*in.AccessType) {
		case "free":
			blog.AccessType = entity.AccessFree
		case "member":
			blog.AccessType = entity.AccessMember
		case "paid_coins":
			blog.AccessType = entity.AccessPaidCoins
		}
	}
	if in.CoinPrice != nil {
		blog.CoinPrice = *in.CoinPrice
	}

	if err := s.blogRepo.Update(ctx, blog); err != nil {
		return nil, apperror.NewInternal(err, "failed to update blog")
	}
	return blog, nil
}

func (s *BlogService) Delete(ctx context.Context, id uint) error {
	blog, err := s.blogRepo.GetByID(ctx, id)
	if err != nil {
		return apperror.NewInternal(err, "failed to fetch blog")
	}
	if blog == nil {
		return apperror.NewNotFound("blog not found")
	}
	return s.blogRepo.Delete(ctx, id)
}

// Analytics

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
	return blog.Views, nil
}

// helpers

func validBlogStatus(s string) bool {
	switch s {
	case "draft", "published", "archived":
		return true
	}
	return false
}

// GetBySlugForReader returns a published blog with access control
func (s *BlogService) GetBySlugForReader(ctx context.Context, slug string, viewerID *uint) (*entity.Blog, error) {
	blog, err := s.blogRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch blog")
	}
	if blog == nil || blog.Status != entity.BlogPublished {
		return nil, apperror.NewNotFound("blog not found")
	}
	return s.applyReaderMask(ctx, blog, viewerID), nil
}

func (s *BlogService) GetByIDForReader(ctx context.Context, id uint, viewerID *uint) (*entity.Blog, error) {
	blog, err := s.blogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch blog")
	}
	if blog == nil || blog.Status != entity.BlogPublished {
		return nil, apperror.NewNotFound("blog not found")
	}
	return s.applyReaderMask(ctx, blog, viewerID), nil
}

func (s *BlogService) applyReaderMask(ctx context.Context, blog *entity.Blog, viewerID *uint) *entity.Blog {
	b := *blog
	if s.canReadFullContent(ctx, &b, viewerID) {
		return &b
	}
	b.Content = ""
	return &b
}

func (s *BlogService) canReadFullContent(ctx context.Context, b *entity.Blog, viewerID *uint) bool {
	switch b.AccessType {
	case entity.AccessFree, "":
		return true
	case entity.AccessMember:
		if viewerID == nil {
			return false
		}
		m, err := s.membershipRepo.GetByUserID(ctx, *viewerID)
		if err != nil || m == nil {
			return false
		}
		return m.IsActive() && m.Plan != entity.PlanFree
	case entity.AccessPaidCoins:
		if viewerID == nil {
			return false
		}
		ok, err := s.coinRepo.HasBlogUnlock(ctx, *viewerID, b.ID)
		return err == nil && ok
	default:
		return false
	}
}
