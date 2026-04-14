// internal/interfaces/http/handler/product_handler.go
package handler

import (
	"net/http"

	"github.com/cureerel/cserver/internal/application/service"
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// POST /api/products — admin
func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := h.productService.Create(c.Request.Context(), service.CreateProductInput{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type, // string — service validates
		Price:       req.Price,
		Currency:    req.Currency,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
	})
	if err != nil {
		respondErr(c, err)
		return
	}
	respondCreated(c, toProductResponse(product))
}

// GET /api/products/:id — public
func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	product, err := h.productService.GetByID(c.Request.Context(), id)
	if err != nil {
		respondErr(c, err)
		return
	}
	respond(c, toProductResponse(product))
}

// GET /api/products — public
func (h *ProductHandler) GetAll(c *gin.Context) {
	page, limit := paginate(c)
	productType := c.Query("type") // optional filter: physical | digital
	products, total, err := h.productService.GetAll(c.Request.Context(), page, limit, productType)
	if err != nil {
		respondErr(c, err)
		return
	}
	list := make([]dto.ProductResponse, len(products))
	for i := range products {
		list[i] = toProductResponse(&products[i])
	}
	c.JSON(http.StatusOK, dto.ProductListResponse{Data: list, Total: total, Page: page, Limit: limit})
}

// PUT /api/products/:id — admin
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := h.productService.Update(c.Request.Context(), id, service.UpdateProductInput{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
		IsActive:    req.IsActive,
	})
	if err != nil {
		respondErr(c, err)
		return
	}
	respond(c, toProductResponse(product))
}

// DELETE /api/products/:id — admin
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	if err := h.productService.Delete(c.Request.Context(), id); err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}

// mapper

func toProductResponse(p *entity.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:          p.ID,
		SKU:         p.SKU,
		Name:        p.Name,
		Description: p.Description,
		Type:        string(p.Type),
		Price:       p.Price,
		PriceUSD:    p.PriceUSD(),
		Currency:    string(p.Currency),
		Stock:       p.Stock,
		ImageURL:    p.ImageURL,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
