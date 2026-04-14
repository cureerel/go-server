// internal/interfaces/dto/blog_dto.go
package dto

type CreateBlogRequest struct {
	Title        string `json:"title"          binding:"required,min=3"`
	Content      string `json:"content"        binding:"required"`
	Keyword      string `json:"keyword"`
	Tag          string `json:"tag"`
	Excerpt      string `json:"excerpt"`
	Thumbnail    string `json:"thumbnail"`
	ThumbnailKey string `json:"thumbnail_key"`
	Status       string `json:"status"`
	AccessType   string `json:"access_type"`
	CoinPrice    int64  `json:"coin_price"`
}

type UpdateBlogRequest struct {
	Title        string `json:"title"`
	Content      string `json:"content"`
	Keyword      string `json:"keyword"`
	Tag          string `json:"tag"`
	Excerpt      string `json:"excerpt"`
	Thumbnail    string `json:"thumbnail"`
	ThumbnailKey string `json:"thumbnail_key"`
	Status       string `json:"status"`
	AccessType   string `json:"access_type"`
	CoinPrice    int64  `json:"coin_price"`
}

type BlogResponse struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	Content    string `json:"content"`
	Keyword    string `json:"keyword"`
	Tag        string `json:"tag"`
	Excerpt    string `json:"excerpt"`
	Thumbnail  string `json:"thumbnail,omitempty"`
	Views      int64  `json:"views"`
	Status     string `json:"status"`
	AccessType string `json:"access_type"`
	CoinPrice  int64  `json:"coin_price"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type BlogListResponse struct {
	Data  []BlogResponse `json:"data"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

type BlogStatsResponse struct {
	BlogID uint  `json:"blog_id"`
	Views  int64 `json:"views"`
}
