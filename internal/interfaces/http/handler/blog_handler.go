package handler

import (
	"net/http"
	"strconv"

	"github.com/cureerel/gotemplate/internal/application/service"
	"github.com/cureerel/gotemplate/internal/domain/entity" 
	"github.com/cureerel/gotemplate/internal/interfaces/dto"

	"github.com/gin-gonic/gin"
)

type BlogHandler struct {
	blogService *service.BlogService
}

func NewBlogHandler(blogService *service.BlogService) *BlogHandler {
	return &BlogHandler{blogService: blogService}
}

// POST /api/blog/create
func (h *BlogHandler) Create(c *gin.Context) {
	var req dto.CreateBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Get authorID from JWT token in real implementation
	authorID := uint(1) // placeholder

	input := service.CreateBlogInput{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: authorID,
		Tags:     req.Tags,
	}

	blog, err := h.blogService.Create(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toBlogResponse(blog))
}

// GET /api/blog/:id
func (h *BlogHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	blog, err := h.blogService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if blog == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
		return
	}

	c.JSON(http.StatusOK, toBlogResponse(blog))
}

// GET /api/blog/slug/:slug
func (h *BlogHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	blog, err := h.blogService.GetBySlug(slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if blog == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
		return
	}

	c.JSON(http.StatusOK, toBlogResponse(blog))
}

// GET /api/blog/list
func (h *BlogHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	blogs, total, err := h.blogService.GetAll(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.BlogListResponse{
		Data:  toBlogListResponse(blogs),
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

// PUT /api/blog/:id
func (h *BlogHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdateBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := service.UpdateBlogInput{
		ID:      uint(id),
		Title:   req.Title,
		Content: req.Content,
		Status:  req.Status,
		Tags:    req.Tags,
	}

	blog, err := h.blogService.Update(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toBlogResponse(blog))
}

// PATCH /api/blog/:id
func (h *BlogHandler) Patch(c *gin.Context) {
	// In REST, PATCH is partial update - reuse Update logic
	// Gin binding with pointers handles partial updates automatically
	h.Update(c)
}

// DELETE /api/blog/:id
func (h *BlogHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.blogService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "blog deleted successfully"})
}

// Helper functions
func toBlogResponse(blog *entity.Blog) dto.BlogResponse {
	return dto.BlogResponse{
		ID:        blog.ID,
		Title:     blog.Title,
		Slug:      blog.Slug,
		Content:   blog.Content,
		AuthorID:  blog.AuthorID,
		Status:    blog.Status,
		Tags:      blog.Tags,
		CreatedAt: blog.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: blog.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func toBlogListResponse(blogs []entity.Blog) []dto.BlogResponse {
	result := make([]dto.BlogResponse, len(blogs))
	for i, blog := range blogs {
		result[i] = toBlogResponse(&blog)
	}
	return result
}