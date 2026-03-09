// internal/application/service/order_service.go
package service

import (
	"context"
	"fmt"

	"github.com/cureerel/gotemplate/internal/domain/entity"
	"github.com/cureerel/gotemplate/internal/domain/repository"
	"github.com/cureerel/gotemplate/pkg/apperror"
)

type OrderService struct {
	orderRepo   repository.OrderRepository
	serviceRepo repository.ServiceRepository
}

func NewOrderService(orderRepo repository.OrderRepository, serviceRepo repository.ServiceRepository) *OrderService {
	return &OrderService{orderRepo: orderRepo, serviceRepo: serviceRepo}
}

type CreateOrderInput struct {
	UserID      uint
	ServiceID   uint   // the service being purchased
	CouponID    *uint
	AffiliateID *uint
	Provider    string // "stripe" | "razorpay"
}

// Create validates the service is live, builds the order with a single item,
// calculates total, and persists. Returns the order ready for payment initiation.
func (s *OrderService) Create(ctx context.Context, in CreateOrderInput) (*entity.Order, error) {
	svc, err := s.serviceRepo.GetByID(ctx, in.ServiceID)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch service")
	}
	if svc == nil {
		return nil, apperror.NewNotFound("service not found")
	}
	if !svc.IsLive() {
		return nil, apperror.NewBadRequest(fmt.Sprintf("service is not available for purchase (status: %s)", svc.Status))
	}

	serviceID := in.ServiceID
	order := &entity.Order{
		UserID:          in.UserID,
		ServiceID:       &serviceID,
		Status:          entity.OrderPending,
		Currency:        "USD",
		PaymentProvider: in.Provider,
		CouponID:        in.CouponID,
		AffiliateID:     in.AffiliateID,
		Items: []entity.OrderItem{
			{
				ServiceID: &serviceID,
				Title:     svc.Title,
				Quantity:  1,
				UnitPrice: svc.PriceUSDCents,
			},
		},
	}
	order.TotalCents = order.CalculateTotal()

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, apperror.NewInternal(err, "failed to create order")
	}
	return order, nil
}

func (s *OrderService) GetByID(ctx context.Context, id uint) (*entity.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch order")
	}
	if order == nil {
		return nil, apperror.NewNotFound("order not found")
	}
	return order, nil
}

// GetMyOrders returns paginated orders for a specific user.
func (s *OrderService) GetMyOrders(ctx context.Context, userID uint, page, limit int) ([]entity.Order, int64, error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 10 }
	return s.orderRepo.GetAll(ctx, repository.OrderFilter{
		Page: page, Limit: limit, UserID: &userID,
	})
}

// GetAll returns all orders (admin use).
func (s *OrderService) GetAll(ctx context.Context, page, limit int, status string) ([]entity.Order, int64, error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 10 }
	return s.orderRepo.GetAll(ctx, repository.OrderFilter{
		Page: page, Limit: limit, Status: status,
	})
}

// UpdateStatus transitions an order to a new status. Admin only.
func (s *OrderService) UpdateStatus(ctx context.Context, id uint, status entity.OrderStatus) error {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil || order == nil {
		return apperror.NewNotFound("order not found")
	}
	if !validOrderStatus(status) {
		return apperror.NewBadRequest("invalid order status")
	}
	return s.orderRepo.UpdateStatus(ctx, id, status)
}

func validOrderStatus(s entity.OrderStatus) bool {
	switch s {
	case entity.OrderPending, entity.OrderConfirmed, entity.OrderCancelled, entity.OrderCompleted:
		return true
	}
	return false
}