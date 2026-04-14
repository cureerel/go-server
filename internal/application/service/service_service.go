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

// ── Input types ───────────────────────────────────────────────

type CreateServiceInput struct {
	OwnerID       uint
	Title         string
	Description   string
	PriceUSDCents int64
	CoverImageURL string
	CoverImageKey string
}

type UpdateServiceInput struct {
	Title         string
	Description   string
	PriceUSDCents int64
	CoverImageURL string
	CoverImageKey string
}

// ── CRUD ──────────────────────────────────────────────────────

// Create submits a new service with status=pending.
// Admin must approve before it becomes live.
func (s *ServiceService) Create(ctx context.Context, in CreateServiceInput) (*entity.Service, error) {
	svc := &entity.Service{
		OwnerID:       in.OwnerID,
		Title:         in.Title,
		Description:   in.Description,
		PriceUSDCents: in.PriceUSDCents,
		Status:        entity.ServiceStatusPending,
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

// GetAll returns services filtered by status (default "live" for public).
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


// ── Analytics ─────────────────────────────────────────────────

func (s *ServiceService) RecordView(ctx context.Context, serviceID uint, ip, ua string) error {
	h := sha256.Sum256([]byte(fmt.Sprintf("%s|%s", ip, ua)))
	_, err := s.repo.RecordView(ctx, serviceID, fmt.Sprintf("%x", h))
	return err
}