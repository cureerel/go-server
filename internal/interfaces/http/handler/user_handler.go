package handler

import (
    "net/http"

    "github.com/cureerel/gotemplate/internal/application/service"
    "github.com/cureerel/gotemplate/internal/interfaces/dto"
    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var req dto.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.userService.CreateUser(req.Username, req.Email)
    if err != nil {
        // Handle custom errors properly here
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, dto.ToUserResponse(user))
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
    users, err := h.userService.GetAllUsers()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    response := make([]dto.UserResponse, len(users))
    for i, u := range users {
        response[i] = dto.ToUserResponse(u)
    }

    c.JSON(http.StatusOK, response)
}