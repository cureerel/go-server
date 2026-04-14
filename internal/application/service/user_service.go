package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/domain/repository"
	"github.com/cureerel/cserver/pkg/apperror"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// GetAll returns paginated users
func (s *UserService) GetAll(ctx context.Context, page, limit int) ([]entity.User, int64, error) {
	users, total, err := s.userRepo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, 0, apperror.NewInternal(err, "failed to fetch users")
	}
	return users, total, nil
}

// Create creates a new user with hashed password
func (s *UserService) Create(ctx context.Context, username, email, password string) (*entity.User, error) {
	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to check existing email")
	}
	if existing != nil {
		return nil, apperror.NewBadRequest("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to hash password")
	}

	user := &entity.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, apperror.NewInternal(err, "failed to create user")
	}
	return user, nil
}

// GetByID fetches a user by ID
func (s *UserService) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch user")
	}
	if user == nil {
		return nil, apperror.NewNotFound("user not found")
	}
	return user, nil
}

// Update modifies a user's details; supports optional fields
func (s *UserService) Update(ctx context.Context, id uint, username, email *string) (*entity.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch user")
	}
	if user == nil {
		return nil, apperror.NewNotFound("user not found")
	}

	if username != nil {
		user.Username = *username
	}
	if email != nil {
		user.Email = *email
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, apperror.NewInternal(err, "failed to update user")
	}
	return user, nil
}

// Delete removes a user by ID
func (s *UserService) Delete(ctx context.Context, id uint) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		return apperror.NewInternal(err, "failed to delete user")
	}
	return nil
}

// GetByEmail fetches a user by email
func (s *UserService) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch user by email")
	}
	return user, nil
}

// Optional: VerifyPassword compares hashed password
func (s *UserService) VerifyPassword(user *entity.User, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return apperror.NewBadRequest("invalid credentials")
	}
	return nil
}
