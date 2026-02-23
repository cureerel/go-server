package dto

type CreateProductRequest struct {
    Name        string `json:"name" binding:"required,min=2,max=200"`
    Description string `json:"description"`
    Type        string `json:"type" binding:"required,oneof=physical digital"`
    Price       int64  `json:"price" binding:"required,gt=0"`
    Currency    string `json:"currency" binding:"required"`
}

type UpdateProductRequest struct {
    Name        *string `json:"name" binding:"omitempty,min=2,max=200"`
    Description *string `json:"description"`
    Price       *int64  `json:"price" binding:"omitempty,gt=0"`
    Currency    *string `json:"currency"`
    IsActive    *bool   `json:"is_active"`
}

type ProductResponse struct {
    ID          uint   `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Type        string `json:"type"`
    Price       int64  `json:"price"`
    Currency    string `json:"currency"`
    IsActive    bool   `json:"is_active"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}