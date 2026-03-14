// internal/application/service/service_service.go
package service

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/cureerel/gotemplate/internal/domain/entity"
	"github.com/cureerel/gotemplate/internal/domain/repository"
	"github.com/cureerel/gotemplate/pkg/apperror"
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

// GetByOwner returns all services for a partner (any status — used on dashboard).
func (s *ServiceService) GetByOwner(ctx context.Context, ownerID uint, page, limit int) ([]entity.Service, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.repo.GetAll(ctx, repository.ServiceFilter{
		Page:    page,
		Limit:   limit,
		OwnerID: &ownerID,
	})
}

// Update applies changes. Only the owner may update. Resetting status to
// pending after an edit is a business choice — currently we keep the existing
// status so approved/live services stay live after minor edits.
func (s *ServiceService) Update(ctx context.Context, id, ownerID uint, in UpdateServiceInput) (*entity.Service, error) {
	svc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch service")
	}
	if svc == nil {
		return nil, apperror.NewNotFound("service not found")
	}
	if svc.OwnerID != ownerID {
		return nil, apperror.NewForbidden("you don't own this service")
	}

	svc.Title = in.Title
	svc.Description = in.Description
	svc.PriceUSDCents = in.PriceUSDCents
	if in.CoverImageURL != "" {
		svc.CoverImageURL = in.CoverImageURL
	}
	if in.CoverImageKey != "" {
		svc.CoverImageKey = in.CoverImageKey
	}

	if err := s.repo.Update(ctx, svc); err != nil {
		return nil, apperror.NewInternal(err, "failed to update service")
	}
	return svc, nil
}

// Delete removes a service. Owner or admin may delete.
func (s *ServiceService) Delete(ctx context.Context, id, callerID uint, callerRole string) error {
	svc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperror.NewInternal(err, "failed to fetch service")
	}
	if svc == nil {
		return apperror.NewNotFound("service not found")
	}
	caller := &entity.User{ID: callerID, Role: callerRole}
	if svc.OwnerID != callerID && !caller.HasRole(entity.RoleAdmin) {
		return apperror.NewForbidden("you don't own this service")
	}
	return s.repo.Delete(ctx, id)
}

// ── Admin approval flow ───────────────────────────────────────

// Approve sets status → approved. Partner must then call SetLive to publish.
func (s *ServiceService) Approve(ctx context.Context, id uint) error {
	svc, err := s.repo.GetByID(ctx, id)
	if err != nil || svc == nil {
		return apperror.NewNotFound("service not found")
	}
	if svc.Status != entity.ServiceStatusPending {
		return apperror.NewBadRequest("only pending services can be approved")
	}
	return s.repo.UpdateStatus(ctx, id, entity.ServiceStatusApproved)
}

// Reject sets status → rejected with a reason (logged via future audit trail).
func (s *ServiceService) Reject(ctx context.Context, id uint) error {
	svc, err := s.repo.GetByID(ctx, id)
	if err != nil || svc == nil {
		return apperror.NewNotFound("service not found")
	}
	if svc.Status != entity.ServiceStatusPending {
		return apperror.NewBadRequest("only pending services can be rejected")
	}
	return s.repo.UpdateStatus(ctx, id, entity.ServiceStatusRejected)
}

// SetLive publishes an approved service. Called by the partner.
func (s *ServiceService) SetLive(ctx context.Context, id, ownerID uint) error {
	svc, err := s.repo.GetByID(ctx, id)
	if err != nil || svc == nil {
		return apperror.NewNotFound("service not found")
	}
	if svc.OwnerID != ownerID {
		return apperror.NewForbidden("you don't own this service")
	}
	if svc.Status != entity.ServiceStatusApproved {
		return apperror.NewBadRequest("service must be approved before going live")
	}
	return s.repo.UpdateStatus(ctx, id, entity.ServiceStatusLive)
}

// Pause takes a live service offline. Owner or admin.
func (s *ServiceService) Pause(ctx context.Context, id, callerID uint, callerRole string) error {
	svc, err := s.repo.GetByID(ctx, id)
	if err != nil || svc == nil {
		return apperror.NewNotFound("service not found")
	}
	caller := &entity.User{ID: callerID, Role: callerRole}
	if svc.OwnerID != callerID && !caller.HasRole(entity.RoleAdmin) {
		return apperror.NewForbidden("you don't own this service")
	}
	if svc.Status != entity.ServiceStatusLive {
		return apperror.NewBadRequest("only live services can be paused")
	}
	return s.repo.UpdateStatus(ctx, id, entity.ServiceStatusPaused)
}

// ── Analytics ─────────────────────────────────────────────────

func (s *ServiceService) RecordView(ctx context.Context, serviceID uint, ip, ua string) error {
	h := sha256.Sum256([]byte(fmt.Sprintf("%s|%s", ip, ua)))
	_, err := s.repo.RecordView(ctx, serviceID, fmt.Sprintf("%x", h))
	return err
}