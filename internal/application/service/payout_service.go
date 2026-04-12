// internal/application/service/payout_service.go
package service

import (
	"context"

	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/domain/repository"
	"github.com/cureerel/cserver/pkg/apperror"
)

type PayoutService struct {
	payoutRepo repository.PayoutRepository
}

func NewPayoutService(payoutRepo repository.PayoutRepository) *PayoutService {
	return &PayoutService{payoutRepo: payoutRepo}
}

func (s *PayoutService) GetByID(ctx context.Context, id uint) (*entity.Payout, error) {
	p, err := s.payoutRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch payout")
	}
	if p == nil {
		return nil, apperror.NewNotFound("payout not found")
	}
	return p, nil
}

func (s *PayoutService) GetAll(ctx context.Context, page, limit int, status string) ([]entity.Payout, int64, error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 10 }
	return s.payoutRepo.GetAll(ctx, page, limit, status)
}

func (s *PayoutService) GetMine(ctx context.Context, recipientID uint, page, limit int) ([]entity.Payout, int64, error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 10 }
	return s.payoutRepo.GetByRecipient(ctx, recipientID, page, limit)
}

// MarkPaid — admin marks a payout as paid with a reference (e.g. bank txn ID).
func (s *PayoutService) MarkPaid(ctx context.Context, id uint, adminID uint, reference string) error {
	p, err := s.payoutRepo.GetByID(ctx, id)
	if err != nil || p == nil {
		return apperror.NewNotFound("payout not found")
	}
	if p.Status == entity.PayoutStatusPaid {
		return nil // idempotent
	}
	if p.Status != entity.PayoutStatusPending {
		return apperror.NewBadRequest("only pending payouts can be marked paid")
	}
	return s.payoutRepo.MarkPaid(ctx, id, adminID, reference)
}

// CreatePartnerPayout queues a partner sale payout when an order completes.
func (s *PayoutService) CreatePartnerPayout(ctx context.Context, recipientID uint, amountCents int64, orderID uint) error {
	if amountCents <= 0 {
		return nil
	}
	oid := orderID
	return s.payoutRepo.Create(ctx, &entity.Payout{
		RecipientID: recipientID,
		Type:        entity.PayoutTypePartnerSale,
		AmountCents: amountCents,
		Status:      entity.PayoutStatusPending,
		OrderID:     &oid,
	})
}