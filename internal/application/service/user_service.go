package service

import (
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

func (s *UserService) GetAll() ([]entity.User, error) {
	return s.userRepo.GetAll()
}

func (s *UserService) Create(name, email, password string) (*entity.User, error) {
	// Check if email exists
	existing, _ := s.userRepo.GetByEmail(email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	user := &entity.User{
		Name:     name,
		Email:    email,
		Password: password, // TODO: Hash password
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetByID(id uint) (*entity.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) Update(id uint, name, email string) (*entity.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		user.Name = name
	}
	if email != "" {
		user.Email = email
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}