package handler

import (
    "net/http"
    "strconv"

    "github.com/cureerel/gotemplate/internal/application/service"
    "github.com/cureerel/gotemplate/internal/domain/entity"
    "github.com/gin-gonic/gin"
)

type ProductHandler struct {
    productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
    return &ProductHandler{productService: productService}
}

type createProductRequest struct {
    Name        string             `json:"name" binding:"required"`
    Description string             `json:"description"`
    Type        entity.ProductType `json:"type" binding:"required"`
    Price       int64              `json:"price" binding:"required,gt=0"`
    Currency    entity.Currency    `json:"currency" binding:"required"`
}

type updateProductRequest struct {
    Name        *string          `json:"name"`
    Description *string          `json:"description"`
    Price       *int64           `json:"price"`
    Currency    *entity.Currency `json:"currency"`
    IsActive    *bool            `json:"is_active"`
}

func (h *ProductHandler) Create(c *gin.Context) {
    var req createProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    product, err := h.productService.Create(c.Request.Context(), service.CreateProductInput{
        Name:        req.Name,
        Description: req.Description,
        Type:        req.Type,
        Price:       req.Price,
        Currency:    req.Currency,
    })
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"data": product})
}

func (h *ProductHandler) GetByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    product, err := h.productService.GetByID(c.Request.Context(), uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": product})
}

func (h *ProductHandler) GetAll(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

    products, total, err := h.productService.GetAll(c.Request.Context(), page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": products, "total": total, "page": page, "limit": limit})
}

func (h *ProductHandler) Update(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    var req updateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    product, err := h.productService.Update(c.Request.Context(), service.UpdateProductInput{
        ID:          uint(id),
        Name:        req.Name,
        Description: req.Description,
        Price:       req.Price,
        Currency:    req.Currency,
        IsActive:    req.IsActive,
    })
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": product})
}

func (h *ProductHandler) Delete(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    if err := h.productService.Delete(c.Request.Context(), uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}