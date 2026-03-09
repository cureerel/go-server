// internal/interfaces/dto/blog_dto.go
package dto

// ── Requests ──────────────────────────────────────────────────

type CreateBlogRequest struct {
	Title         string `json:"title"          binding:"required,min=2,max=200"`
	Content       string `json:"content"        binding:"required"`
	Tags          string `json:"tags"`
	CoverImageURL string `json:"cover_image_url"`
	CoverImageKey string `json:"cover_image_key"`
}

type UpdateBlogRequest struct {
	Title         *string `json:"title"`
	Content       *string `json:"content"`
	Status        *string `json:"status"`
	Tags          *string `json:"tags"`
	CoverImageURL *string `json:"cover_image_url"`
	CoverImageKey *string `json:"cover_image_key"`
}

// ── Responses ─────────────────────────────────────────────────

type BlogResponse struct {
	ID            uint   `json:"id"`
	Title         string `json:"title"`
	Slug          string `json:"slug"`
	Content       string `json:"content"`
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

// ── Upload ────────────────────────────────────────────────────

type UploadResponse struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

type DeleteUploadRequest struct {
	Key string `json:"key" binding:"required"`
}