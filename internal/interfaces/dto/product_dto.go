// internal/interfaces/dto/product_dto.go
package dto

// Requests

type CreateProductRequest struct {
	Name        string `json:"name"        binding:"required,min=2,max=200"`
	Description string `json:"description"`
	Type        string `json:"type"        binding:"required,oneof=physical digital"`
	Price       int64  `json:"price"       binding:"required,gt=0"`
	Currency    string `json:"currency"` // default USD
	Stock       int    `json:"stock"`    // ignored for digital
	ImageURL    string `json:"image_url"`
}

type UpdateProductRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Price       *int64  `json:"price"`
	Stock       *int    `json:"stock"`
	ImageURL    *string `json:"image_url"`
	IsActive    *bool   `json:"is_active"`
}

// Responses

type ProductResponse struct {
	ID          uint    `json:"id"`
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	Price       int64   `json:"price"`
	PriceUSD    float64 `json:"price_usd"`
	Currency    string  `json:"currency"`
	Stock       int     `json:"stock"`
	ImageURL    string  `json:"image_url,omitempty"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type ProductListResponse struct {
	Data  []ProductResponse `json:"data"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
}
