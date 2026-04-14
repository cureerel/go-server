// internal/interfaces/http/handler/blog_handler.go
package handler

import (
	"net/http"

	"github.com/cureerel/cserver/internal/application/service"
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type BlogHandler struct {
	svc  *service.BlogService
	coin *service.CoinService
}

func NewBlogHandler(svc *service.BlogService, coin *service.CoinService) *BlogHandler {
	return &BlogHandler{svc: svc, coin: coin}
}

// POST /api/blogs — admin only
func (h *BlogHandler) Create(c *gin.Context) {
	var req dto.CreateBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	blog, err := h.svc.Create(c.Request.Context(), service.CreateBlogInput{
		Title:        req.Title,
		Content:      req.Content,
		Keyword:      req.Keyword,
		Tag:          req.Tag,
		Excerpt:      req.Excerpt,
		Thumbnail:    req.Thumbnail,
		ThumbnailKey: req.ThumbnailKey,
		Status:       req.Status,
		AccessType:   req.AccessType,
		CoinPrice:    req.CoinPrice,
	})
	if err != nil {
		respondErr(c, err)
		return
	}
	respondCreated(c, toBlogResponse(blog))
}

// GET /api/blogs — public paginated
func (h *BlogHandler) GetAll(c *gin.Context) {
	page, limit := paginate(c)
	search := c.Query("search")
	tag := c.Query("tag")
	blogs, total, err := h.svc.GetAll(c.Request.Context(), page, limit, search, tag)
	if err != nil {
		respondErr(c, err)
		return
	}
	list := make([]dto.BlogResponse, len(blogs))
	for i := range blogs {
		list[i] = toBlogResponse(&blogs[i])
	}
	c.JSON(http.StatusOK, dto.BlogListResponse{Data: list, Total: total, Page: page, Limit: limit})
}

// GET /api/blogs/:id — public (published only; optional auth unlocks member/paid)
func (h *BlogHandler) GetByID(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	var viewer *uint
	if uid, ok := getUID(c); ok {
		viewer = &uid
	}
	blog, err := h.svc.GetByIDForReader(c.Request.Context(), id, viewer)
	if err != nil {
		respondErr(c, err)
		return
	}
	go func() { _ = h.svc.RecordView(c.Request.Context(), id, c.ClientIP(), c.Request.UserAgent()) }()
	respond(c, toBlogResponse(blog))
}

// GET /api/blogs/slug/:slug — public
func (h *BlogHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	var viewer *uint
	if uid, ok := getUID(c); ok {
		viewer = &uid
	}
	blog, err := h.svc.GetBySlugForReader(c.Request.Context(), slug, viewer)
	if err != nil {
		respondErr(c, err)
		return
	}
	go func() { _ = h.svc.RecordView(c.Request.Context(), blog.ID, c.ClientIP(), c.Request.UserAgent()) }()
	respond(c, toBlogResponse(blog))
}

// GET /api/blogs/mine — admin (own blogs, any status)
func (h *BlogHandler) GetMine(c *gin.Context) {
	uid, ok := getUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	page, limit := paginate(c)
	blogs, total, err := h.svc.GetMine(c.Request.Context(), uid, page, limit)
	if err != nil {
		respondErr(c, err)
		return
	}
	list := make([]dto.BlogResponse, len(blogs))
	for i := range blogs {
		list[i] = toBlogResponse(&blogs[i])
	}
	c.JSON(http.StatusOK, dto.BlogListResponse{Data: list, Total: total, Page: page, Limit: limit})
}

// PUT /api/blogs/:id — admin only
func (h *BlogHandler) Update(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	var req dto.UpdateBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, _ := getUID(c)
	blog, err := h.svc.Update(c.Request.Context(), service.UpdateBlogInput{
		ID:           id,
		CallerID:     uid,
		CallerRole:   getRole(c),
		Title:        &req.Title,
		Content:      &req.Content,
		Keyword:      &req.Keyword,
		Tag:          &req.Tag,
		Excerpt:      &req.Excerpt,
		Thumbnail:    &req.Thumbnail,
		ThumbnailKey: &req.ThumbnailKey,
		Status:       &req.Status,
		AccessType:   &req.AccessType,
		CoinPrice:    &req.CoinPrice,
	})
	if err != nil {
		respondErr(c, err)
		return
	}
	respond(c, toBlogResponse(blog))
}

// PATCH /api/blogs/:id — alias
func (h *BlogHandler) Patch(c *gin.Context) { h.Update(c) }

// DELETE /api/blogs/:id — admin only
func (h *BlogHandler) Delete(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "blog deleted"})
}

// GET /api/blogs/:id/stats — admin only
func (h *BlogHandler) GetStats(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	views, err := h.svc.GetStats(c.Request.Context(), id)
	if err != nil {
		respondErr(c, err)
		return
	}
	respond(c, dto.BlogStatsResponse{BlogID: id, Views: views})
}

// mapper

func toBlogResponse(b *entity.Blog) dto.BlogResponse {
	return dto.BlogResponse{
		ID:         b.ID,
		Title:      b.Title,
		Slug:       b.Slug,
		Content:    b.Content,
		Keyword:    b.Keyword,
		Tag:        b.Tag,
		Excerpt:    b.Excerpt,
		Thumbnail:  b.Thumbnail,
		Views:      b.Views,
		Status:     string(b.Status),
		AccessType: string(b.AccessType),
		CoinPrice:  b.CoinPrice,
		CreatedAt:  b.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:  b.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// POST /api/blogs/:id/unlock — spend coins for paid_coins posts
func (h *BlogHandler) UnlockPaidBlog(c *gin.Context) {
	if h.coin == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "coins unavailable"})
		return
	}
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	uid, ok := getUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	blog, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		respondErr(c, err)
		return
	}
	if blog == nil || blog.Status != entity.BlogPublished || blog.AccessType != entity.AccessPaidCoins {
		c.JSON(http.StatusBadRequest, gin.H{"error": "post is not unlockable"})
		return
	}
	if err := h.coin.UnlockBlog(c.Request.Context(), uid, blog.ID, blog.CoinPrice); err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "unlocked"})
}
