package handler

import (
    "fmt"
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

func (h *BlogHandler) Create(c *gin.Context) {
    var req dto.CreateBlogRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    uid, _ := strconv.ParseUint(fmt.Sprintf("%v", userID), 10, 32)

    blog, err := h.blogService.Create(c.Request.Context(), service.CreateBlogInput{
        Title:    req.Title,
        Content:  req.Content,
        AuthorID: uint(uid),
        Tags:     req.Tags,
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"data": toBlogResponse(blog)})
}

func (h *BlogHandler) GetByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    blog, err := h.blogService.GetByID(c.Request.Context(), uint(id))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if blog == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": toBlogResponse(blog)})
}

func (h *BlogHandler) GetBySlug(c *gin.Context) {
    slug := c.Param("slug")
    blog, err := h.blogService.GetBySlug(c.Request.Context(), slug)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if blog == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": toBlogResponse(blog)})
}

func (h *BlogHandler) GetAll(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

    blogs, total, err := h.blogService.GetAll(c.Request.Context(), page, limit)
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

    blog, err := h.blogService.Update(c.Request.Context(), service.UpdateBlogInput{
        ID:      uint(id),
        Title:   req.Title,
        Content: req.Content,
        Status:  req.Status,
        Tags:    req.Tags,
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": toBlogResponse(blog)})
}

func (h *BlogHandler) Patch(c *gin.Context) {
    h.Update(c)
}

func (h *BlogHandler) Delete(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    if err := h.blogService.Delete(c.Request.Context(), uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "blog deleted successfully"})
}

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