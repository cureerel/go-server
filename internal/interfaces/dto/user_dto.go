package dto

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" binding:"omitempty,min=2,max=100"`
	Email string `json:"email" binding:"omitempty,email"`
}

type UserResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	IsActive   bool   `json:"is_active"`
	IsVerified bool   `json:"is_verified"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

