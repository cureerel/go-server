// internal/interfaces/dto/service_dto.go
package dto

// Requests

type CreateServiceRequest struct {
	Title         string `json:"title"           binding:"required,min=3"`
	Description   string `json:"description"     binding:"required"`
	PriceUSDCents int64  `json:"price_usd_cents" binding:"required,min=1"`
	CoverImageURL string `json:"cover_image_url"`
	CoverImageKey string `json:"cover_image_key"`
}

type UpdateServiceRequest struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	PriceUSDCents int64  `json:"price_usd_cents"`
	CoverImageURL string `json:"cover_image_url"`
	CoverImageKey string `json:"cover_image_key"`
}

// Responses

type ServiceResponse struct {
	ID            uint    `json:"id"`
	OwnerID       uint    `json:"owner_id"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	PriceUSDCents int64   `json:"price_usd_cents"`
	PriceUSD      float64 `json:"price_usd"`
	Status        string  `json:"status"`
	CoverImageURL string  `json:"cover_image_url,omitempty"`
	ViewsTotal    int64   `json:"views_total"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type ServiceListResponse struct {
	Data  []ServiceResponse `json:"data"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
}
