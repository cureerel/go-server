// internal/application/service/order_service.go
package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/domain/repository"
	"github.com/cureerel/cserver/pkg/apperror"
	"gorm.io/gorm"
)

type OrderService struct {
	db          *gorm.DB
	orderRepo   repository.OrderRepository
	serviceRepo repository.ServiceRepository
	couponRepo  repository.CouponRepository
	coinRepo    repository.CoinRepository
}

func NewOrderService(
	db *gorm.DB,
	orderRepo repository.OrderRepository,
	serviceRepo repository.ServiceRepository,
	couponRepo repository.CouponRepository,
	coinRepo repository.CoinRepository,
) *OrderService {
	return &OrderService{
		db:          db,
		orderRepo:   orderRepo,
		serviceRepo: serviceRepo,
		couponRepo:  couponRepo,
		coinRepo:    coinRepo,
	}
}

type CreateOrderInput struct {
	UserID      uint
	ServiceID   uint
	CouponCode  string
	AffiliateID *uint
}

func (s *OrderService) Create(ctx context.Context, in CreateOrderInput) (*entity.Order, error) {
	svc, err := s.serviceRepo.GetByID(ctx, in.ServiceID)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch service")
	}
	if svc == nil {
		return nil, apperror.NewNotFound("service not found")
	}
	if !svc.IsLive() {
		return nil, apperror.NewBadRequest(fmt.Sprintf("service is not available (status: %s)", svc.Status))
	}

	serviceID := in.ServiceID
	totalCents := svc.PriceUSDCents

	var couponID *uint
	if in.CouponCode != "" {
		c, err := s.couponRepo.GetByCode(ctx, in.CouponCode)
		if err != nil {
			return nil, apperror.NewInternal(err, "failed to look up coupon")
		}
		if c != nil && c.IsValid() {
			discount := c.ApplyDiscount(totalCents)
			totalCents -= discount
			if totalCents < 0 {
				totalCents = 0
			}
			couponID = &c.ID
		}
	}

	coinCost := totalCents
	if coinCost <= 0 {
		return nil, apperror.NewBadRequest("nothing to pay — order total is zero")
	}

	order := &entity.Order{
		UserID:          in.UserID,
		ServiceID:       &serviceID,
		Status:          entity.OrderCompleted,
		TotalCents:      totalCents,
		Currency:        "USD",
		PaymentProvider: string(entity.ProviderCoins),
		CouponID:        couponID,
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

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		refID := serviceID
		_, err := s.coinRepo.Debit(ctx, tx, in.UserID, coinCost, "service_purchase", "service", &refID)
		if err != nil {
			return err
		}
		return s.orderRepo.CreateWithTx(ctx, tx, order)
	})
	if err != nil {
		if errors.Is(err, repository.ErrInsufficientCoins) {
			return nil, apperror.NewBadRequest("insufficient coins — top up your wallet first")
		}
		return nil, apperror.NewInternal(err, "failed to place order")
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

func (s *OrderService) GetMyOrders(ctx context.Context, userID uint, page, limit int) ([]entity.Order, int64, error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 10 }
	return s.orderRepo.GetAll(ctx, repository.OrderFilter{
		Page: page, Limit: limit, UserID: &userID,
	})
}

func (s *OrderService) GetAll(ctx context.Context, page, limit int, status string) ([]entity.Order, int64, error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 10 }
	return s.orderRepo.GetAll(ctx, repository.OrderFilter{
		Page: page, Limit: limit, Status: status,
	})
}

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