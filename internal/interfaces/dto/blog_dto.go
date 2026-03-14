// internal/interfaces/dto/blog_dto.go
package dto

type CreateBlogRequest struct {
	Title         string `json:"title"          binding:"required,min=3"`
	Content       string `json:"content"        binding:"required"`
	Excerpt       string `json:"excerpt"`        
	Tags          string `json:"tags"`
	CoverImageURL string `json:"cover_image_url"`
	CoverImageKey string `json:"cover_image_key"`
	Status        string `json:"status"`         
}

type UpdateBlogRequest struct {
	Title         string `json:"title"`
	Content       string `json:"content"`
	Excerpt       string `json:"excerpt"`
	Tags          string `json:"tags"`
	CoverImageURL string `json:"cover_image_url"`
	CoverImageKey string `json:"cover_image_key"`
	Status        string `json:"status"`
}

type BlogResponse struct {
	ID            uint   `json:"id"`
	Title         string `json:"title"`
	Slug          string `json:"slug"`
	Content       string `json:"content"`
	Excerpt       string `json:"excerpt"`
	AuthorID      uint   `json:"author_id"`
	Status        string `json:"status"`
	Tags          string `json:"tags"`
	CoverImageURL string `json:"cover_image_url,omitempty"`
	ViewsTotal    int64  `json:"views_total"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type BlogListResponse struct {
	Data  []BlogResponse `json:"data"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

type BlogStatsResponse struct {
	BlogID     uint  `json:"blog_id"`
	ViewsTotal int64 `json:"views_total"`
}