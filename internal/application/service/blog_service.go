// internal/application/service/blog_service.go
package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/domain/repository"
	"github.com/cureerel/cserver/pkg/apperror"
	"github.com/cureerel/cserver/pkg/utils"
)

type BlogService struct {
	blogRepo         repository.BlogRepository
	coinRepo         repository.CoinRepository
	membershipRepo   repository.MembershipRepository
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

// ── Input types ───────────────────────────────────────────────

type CreateBlogInput struct {
	Title         string
	Content       string
	Excerpt       string
	Status        string
	AuthorID      uint
	Tags          string
	CoverImageURL string
	CoverImageKey string
	AccessType    string
	CoinPrice     int64
	CoAuthorIDs   []uint
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

	// Publishing goes through review — new posts start as draft only.
	status := entity.BlogDraft
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
		Title:         in.Title,
		Slug:          slug,
		Content:       in.Content,
		Excerpt:       in.Excerpt,
		AuthorID:      in.AuthorID,
		Tags:          in.Tags,
		Status:        status,
		AccessType:    at,
		CoinPrice:     coinPrice,
		CoverImageURL: in.CoverImageURL,
		CoverImageKey: in.CoverImageKey,
	}
	if err := s.blogRepo.Create(ctx, blog); err != nil {
		return nil, apperror.NewInternal(err, "failed to create blog")
	}
	for _, uid := range in.CoAuthorIDs {
		if uid == 0 || uid == in.AuthorID {
			continue
		}
		_ = s.blogRepo.AddCoAuthor(ctx, blog.ID, uid)
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
	can, err := s.blogRepo.UserCanEditBlog(ctx, in.ID, in.CallerID)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to check blog access")
	}
	if !can && !caller.HasRole(entity.RoleAdmin) {
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
			return nil, apperror.NewBadRequest("invalid blog status")
		}
		newSt := entity.BlogStatus(*in.Status)
		if newSt == entity.BlogPublished && !caller.HasRole(entity.RoleAdmin) && !caller.HasRole(entity.RoleSuperAdmin) {
			return nil, apperror.NewForbidden("publishing is done by reviewers after submit-for-review")
		}
		blog.Status = newSt
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
	can, err := s.blogRepo.UserCanEditBlog(ctx, id, callerID)
	if err != nil {
		return apperror.NewInternal(err, "failed to check blog access")
	}
	if !can && !caller.HasRole(entity.RoleAdmin) {
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
	case "draft", "in_review", "published", "rejected", "archived":
		return true
	}
	return false
}

// GetBySlugForReader returns published posts; masks body when access rules block the viewer.
func (s *BlogService) GetBySlugForReader(ctx context.Context, slug string, viewerID *uint) (*entity.Blog, error) {
	blog, err := s.blogRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch blog")
	}
	if blog == nil {
		return nil, apperror.NewNotFound("blog not found")
	}
	if blog.Status != entity.BlogPublished {
		return nil, apperror.NewNotFound("blog not found")
	}
	return s.applyReaderMask(ctx, blog, viewerID), nil
}

func (s *BlogService) GetByIDForReader(ctx context.Context, id uint, viewerID *uint) (*entity.Blog, error) {
	blog, err := s.blogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch blog")
	}
	if blog == nil {
		return nil, apperror.NewNotFound("blog not found")
	}
	if blog.Status != entity.BlogPublished {
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

// SubmitForReview moves draft/rejected → in_review.
func (s *BlogService) SubmitForReview(ctx context.Context, blogID, callerID uint, callerRole string) (*entity.Blog, error) {
	blog, err := s.blogRepo.GetByID(ctx, blogID)
	if err != nil || blog == nil {
		return nil, apperror.NewNotFound("blog not found")
	}
	caller := &entity.User{ID: callerID, Role: callerRole}
	can, err := s.blogRepo.UserCanEditBlog(ctx, blogID, callerID)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to check blog access")
	}
	if !can && !caller.HasRole(entity.RoleAdmin) {
		return nil, apperror.NewForbidden("forbidden")
	}
	if blog.Status != entity.BlogDraft && blog.Status != entity.BlogRejected {
		return nil, apperror.NewBadRequest("only draft or rejected posts can be submitted")
	}
	now := time.Now().UTC()
	blog.Status = entity.BlogInReview
	blog.SubmittedForReviewAt = &now
	if err := s.blogRepo.Update(ctx, blog); err != nil {
		return nil, apperror.NewInternal(err, "failed to update blog")
	}
	return blog, nil
}

// ReviewApprove publishes an in_review post (reviewer/admin/superadmin).
func (s *BlogService) ReviewApprove(ctx context.Context, blogID, reviewerID uint) (*entity.Blog, error) {
	blog, err := s.blogRepo.GetByID(ctx, blogID)
	if err != nil || blog == nil {
		return nil, apperror.NewNotFound("blog not found")
	}
	if blog.Status != entity.BlogInReview {
		return nil, apperror.NewBadRequest("post is not awaiting review")
	}
	now := time.Now().UTC()
	blog.Status = entity.BlogPublished
	blog.PublishedAt = &now
	blog.ReviewedByID = &reviewerID
	blog.ReviewNote = ""
	if err := s.blogRepo.Update(ctx, blog); err != nil {
		return nil, apperror.NewInternal(err, "failed to publish")
	}
	return blog, nil
}

// ReviewReject sends back to author.
func (s *BlogService) ReviewReject(ctx context.Context, blogID, reviewerID uint, note string) (*entity.Blog, error) {
	blog, err := s.blogRepo.GetByID(ctx, blogID)
	if err != nil || blog == nil {
		return nil, apperror.NewNotFound("blog not found")
	}
	if blog.Status != entity.BlogInReview {
		return nil, apperror.NewBadRequest("post is not awaiting review")
	}
	blog.Status = entity.BlogRejected
	blog.ReviewedByID = &reviewerID
	blog.ReviewNote = note
	if err := s.blogRepo.Update(ctx, blog); err != nil {
		return nil, apperror.NewInternal(err, "failed to update blog")
	}
	return blog, nil
}

// ListReviewQueue for reviewers.
func (s *BlogService) ListReviewQueue(ctx context.Context, page, limit int) ([]entity.Blog, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.blogRepo.ListByStatus(ctx, string(entity.BlogInReview), page, limit)
}