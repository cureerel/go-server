package dto

// MUST ADD THIS IMPORT:
import "github.com/cureerel/gotemplate/internal/domain/entity"

type CreateUserRequest struct {
    Username string `json:"username" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
}

type UserResponse struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

func ToUserResponse(u *entity.User) UserResponse {
    return UserResponse{
        ID:       u.ID,
        Username: u.Username,
        Email:    u.Email,
    }
}