// internal/application/service/payment_service.go
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/domain/repository"
	"github.com/cureerel/cserver/pkg/apperror"
)

type PaymentService struct {
	paymentRepo repository.PaymentRepository
	orderRepo   repository.OrderRepository
}

func NewPaymentService(paymentRepo repository.PaymentRepository, orderRepo repository.OrderRepository) *PaymentService {
	return &PaymentService{paymentRepo: paymentRepo, orderRepo: orderRepo}
}

// InitPayment creates a pending payment record for an order.
// Called right after order creation, before the provider is contacted.
func (s *PaymentService) InitPayment(ctx context.Context, order *entity.Order, provider entity.PaymentProvider, customerEmail string) (*entity.Payment, error) {
	// Generate a deterministic payment ID: provider-orderID-timestamp
	id := fmt.Sprintf("%s-%d-%d", string(provider), order.ID, time.Now().UnixMilli())

	payment := &entity.Payment{
		ID:            id,
		OrderID:       order.ID,
		UserID:        order.UserID,
		AmountCents:   order.TotalCents,
		Currency:      order.Currency,
		Status:        entity.PaymentPending,
		Provider:      provider,
		CustomerEmail: customerEmail,
		Description:   fmt.Sprintf("Order #%d", order.ID),
	}
	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, apperror.NewInternal(err, "failed to create payment record")
	}
	// Attach provider to order
	_ = s.orderRepo.AttachPaymentProvider(ctx, order.ID, string(provider))
	return payment, nil
}

func (s *PaymentService) GetByID(ctx context.Context, id string) (*entity.Payment, error) {
	p, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch payment")
	}
	if p == nil {
		return nil, apperror.NewNotFound("payment not found")
	}
	return p, nil
}

func (s *PaymentService) GetByOrderID(ctx context.Context, orderID uint) (*entity.Payment, error) {
	p, err := s.paymentRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch payment")
	}
	if p == nil {
		return nil, apperror.NewNotFound("payment not found for this order")
	}
	return p, nil
}

// GetAll returns all payments (admin use).
func (s *PaymentService) GetAll(ctx context.Context, page, limit int, status string) ([]entity.Payment, int64, error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 10 }
	return s.paymentRepo.GetAll(ctx, repository.PaymentFilter{
		Page: page, Limit: limit, Status: status,
	})
}

// MarkCompleted is called by webhook or admin when provider confirms payment.
// Also transitions the linked order to confirmed.
func (s *PaymentService) MarkCompleted(ctx context.Context, id string) error {
	p, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if p.IsCompleted() {
		return nil // idempotent
	}
	if err := s.paymentRepo.UpdateStatus(ctx, id, entity.PaymentCompleted); err != nil {
		return apperror.NewInternal(err, "failed to update payment status")
	}
	_ = s.orderRepo.UpdateStatus(ctx, p.OrderID, entity.OrderConfirmed)
	return nil
}

// MarkFailed is called when provider reports failure.
func (s *PaymentService) MarkFailed(ctx context.Context, id string) error {
	p, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.paymentRepo.UpdateStatus(ctx, id, entity.PaymentFailed); err != nil {
		return apperror.NewInternal(err, "failed to update payment status")
	}
	_ = s.orderRepo.UpdateStatus(ctx, p.OrderID, entity.OrderCancelled)
	return nil
}

// Refund marks a completed payment as refunded.
func (s *PaymentService) Refund(ctx context.Context, id string, refundID string) error {
	p, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !p.IsCompleted() {
		return apperror.NewBadRequest("only completed payments can be refunded")
	}
	if err := s.paymentRepo.MarkRefunded(ctx, id, refundID); err != nil {
		return apperror.NewInternal(err, "failed to mark refund")
	}
	_ = s.orderRepo.UpdateStatus(ctx, p.OrderID, entity.OrderCancelled)
	return nil
}

// SetProviderTxnID stores the provider's transaction reference after initiation.
func (s *PaymentService) SetProviderTxnID(ctx context.Context, id string, txnID string) error {
	return s.paymentRepo.SetProviderTxnID(ctx, id, txnID)
}