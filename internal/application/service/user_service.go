package service

import (
    "errors" 

    "github.com/cureerel/gotemplate/internal/domain/entity"
    "github.com/cureerel/gotemplate/internal/domain/repository"
    customErrors "github.com/cureerel/gotemplate/pkg/errors" // Alias your custom errors to avoid conflict
)

type UserService struct {
    repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) CreateUser(username, email string) (*entity.User, error) {
    if username == "" {
        // Use standard errors.New for the message, pass it to your custom wrapper
        return nil, customErrors.NewBadRequest(errors.New("username cannot be empty"))
    }

    user := &entity.User{
        Username: username,
        Email:    email,
    }

    if err := s.repo.Create(user); err != nil {
        return nil, customErrors.NewInternal(err)
    }

    return user, nil
}

func (s *UserService) GetAllUsers() ([]*entity.User, error) {
    return s.repo.FindAll()
}