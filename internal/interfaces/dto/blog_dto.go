package dto

type CreateBlogRequest struct {
	Title   string `json:"title" binding:"required,min=3,max=200"`
	Content string `json:"content" binding:"required"`
	Tags    string `json:"tags"`
}

type UpdateBlogRequest struct {
	Title   *string `json:"title" binding:"omitempty,min=3,max=200"`
	Content *string `json:"content"`
	Status  *string `json:"status" binding:"omitempty,oneof=draft published archived"`
	Tags    *string `json:"tags"`
}

type BlogResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	Content   string `json:"content"`
	AuthorID  uint   `json:"author_id"`
	Status    string `json:"status"`
	Tags      string `json:"tags"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type BlogListResponse struct {
	Data  []BlogResponse `json:"data"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}