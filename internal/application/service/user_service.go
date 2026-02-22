package service

import (
    "context"
    "errors"

    "github.com/cureerel/gotemplate/internal/domain/entity"
    "github.com/cureerel/gotemplate/internal/domain/repository"
)

type UserService struct {
    userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
    return &UserService{userRepo: userRepo}
}

func (s *UserService) GetAll(ctx context.Context, page, limit int) ([]entity.User, int64, error) {
    return s.userRepo.GetAll(ctx, page, limit)
}

func (s *UserService) Create(ctx context.Context, name, email, password string) (*entity.User, error) {
    existing, err := s.userRepo.GetByEmail(ctx, email)
    if err != nil {
        return nil, err
    }
    if existing != nil {
        return nil, errors.New("email already exists")
    }

    user := &entity.User{
        Name:     name,
        Email:    email,
        Password: password, // TODO: hash password
    }

    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }
    return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*entity.User, error) {
    return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, id uint, name, email string) (*entity.User, error) {
    user, err := s.userRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("user not found")
    }

    if name != "" {
        user.Name = name
    }
    if email != "" {
        user.Email = email
    }

    if err := s.userRepo.Update(ctx, user); err != nil {
        return nil, err
    }
    return user, nil
}

func (s *UserService) Delete(ctx context.Context, id uint) error {
    return s.userRepo.Delete(ctx, id)
}