// internal/application/service/service_service.go
package service

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/domain/repository"
	"github.com/cureerel/cserver/pkg/apperror"
)

type ServiceService struct {
	repo repository.ServiceRepository
}

func NewServiceService(repo repository.ServiceRepository) *ServiceService {
	return &ServiceService{repo: repo}
}

type CreateServiceInput struct {
	OwnerID       uint
	Title         string
	Description   string
	PriceUSDCents int64
	CoverImageURL string
	CoverImageKey string
}

// Create — admin creates a service, starts as live immediately.
func (s *ServiceService) Create(ctx context.Context, in CreateServiceInput) (*entity.Service, error) {
	svc := &entity.Service{
		OwnerID:       in.OwnerID,
		Title:         in.Title,
		Description:   in.Description,
		PriceUSDCents: in.PriceUSDCents,
		Status:        entity.ServiceStatusLive,
		CoverImageURL: in.CoverImageURL,
		CoverImageKey: in.CoverImageKey,
	}
	if err := s.repo.Create(ctx, svc); err != nil {
		return nil, apperror.NewInternal(err, "failed to create service")
	}
	return svc, nil
}

func (s *ServiceService) GetByID(ctx context.Context, id uint) (*entity.Service, error) {
	svc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch service")
	}
	if svc == nil {
		return nil, apperror.NewNotFound("service not found")
	}
	return svc, nil
}

func (s *ServiceService) GetAll(ctx context.Context, page, limit int, status, search string) ([]entity.Service, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.repo.GetAll(ctx, repository.ServiceFilter{
		Page:   page,
		Limit:  limit,
		Status: status,
		Search: search,
	})
}

// Delete — admin removes a service.
func (s *ServiceService) Delete(ctx context.Context, id, callerID uint, callerRole string) error {
	svc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperror.NewInternal(err, "failed to fetch service")
	}
	if svc == nil {
		return apperror.NewNotFound("service not found")
	}
	caller := &entity.User{ID: callerID, Role: callerRole}
	if !caller.HasRole(entity.RoleAdmin) {
		return apperror.NewForbidden("admin only")
	}
	return s.repo.Delete(ctx, id)
}

// SetLive — admin sets a service live.
func (s *ServiceService) SetLive(ctx context.Context, id, callerID uint) error {
	svc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperror.NewInternal(err, "failed to fetch service")
	}
	if svc == nil {
		return apperror.NewNotFound("service not found")
	}
	return s.repo.UpdateStatus(ctx, id, entity.ServiceStatusLive)
}

// Pause — admin pauses a service.
func (s *ServiceService) Pause(ctx context.Context, id, callerID uint, callerRole string) error {
	svc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperror.NewInternal(err, "failed to fetch service")
	}
	if svc == nil {
		return apperror.NewNotFound("service not found")
	}
	return s.repo.UpdateStatus(ctx, id, entity.ServiceStatusPaused)
}

// Analytics

func (s *ServiceService) RecordView(ctx context.Context, serviceID uint, ip, ua string) error {
	h := sha256.Sum256([]byte(fmt.Sprintf("%s|%s", ip, ua)))
	_, err := s.repo.RecordView(ctx, serviceID, fmt.Sprintf("%x", h))
	return err
}
