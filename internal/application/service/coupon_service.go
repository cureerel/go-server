// internal/application/service/coupon_service.go
package service

import (
	"context"
	"strings"

	"github.com/cureerel/gotemplate/internal/domain/entity"
	"github.com/cureerel/gotemplate/internal/domain/repository"
	"github.com/cureerel/gotemplate/pkg/apperror"
)

type CouponService struct {
	couponRepo      repository.CouponRepository
	couponUsageRepo repository.CouponUsageRepository
	payoutRepo      repository.PayoutRepository
}

func NewCouponService(
	couponRepo repository.CouponRepository,
	couponUsageRepo repository.CouponUsageRepository,
	payoutRepo repository.PayoutRepository,
) *CouponService {
	return &CouponService{
		couponRepo:      couponRepo,
		couponUsageRepo: couponUsageRepo,
		payoutRepo:      payoutRepo,
	}
}

type CreateCouponInput struct {
	CreatorID        uint
	Code             string
	Type             string
	DiscountUSDCents int64
	MaxDiscountCents int64
	CommissionPct    float64
	UsageLimit       *int
	ExpiresAt        *string // RFC3339 optional
}

// Create — any partner+ can create a coupon; it starts as pending.
func (s *CouponService) Create(ctx context.Context, in CreateCouponInput) (*entity.Coupon, error) {
	code := strings.ToUpper(strings.TrimSpace(in.Code))
	if code == "" {
		return nil, apperror.NewBadRequest("coupon code is required")
	}
	existing, err := s.couponRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to check coupon code")
	}
	if existing != nil {
		return nil, apperror.NewConflict("coupon code already exists")
	}
	if !validCouponType(in.Type) {
		return nil, apperror.NewBadRequest("type must be discount, affiliate, or both")
	}

	c := &entity.Coupon{
		CreatorID:        in.CreatorID,
		Code:             code,
		Type:             in.Type,
		DiscountUSDCents: in.DiscountUSDCents,
		MaxDiscountCents: in.MaxDiscountCents,
		CommissionPct:    in.CommissionPct,
		Status:           entity.CouponStatusPending,
		UsageLimit:       in.UsageLimit,
	}
	if err := s.couponRepo.Create(ctx, c); err != nil {
		return nil, apperror.NewInternal(err, "failed to create coupon")
	}
	return c, nil
}

func (s *CouponService) GetByID(ctx context.Context, id uint) (*entity.Coupon, error) {
	c, err := s.couponRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch coupon")
	}
	if c == nil {
		return nil, apperror.NewNotFound("coupon not found")
	}
	return c, nil
}

// Validate returns the coupon if it is valid for use. Used at checkout.
func (s *CouponService) Validate(ctx context.Context, code string) (*entity.Coupon, error) {
	c, err := s.couponRepo.GetByCode(ctx, strings.ToUpper(strings.TrimSpace(code)))
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to look up coupon")
	}
	if c == nil || !c.IsValid() {
		return nil, apperror.NewBadRequest("invalid or expired coupon code")
	}
	return c, nil
}

func (s *CouponService) GetAll(ctx context.Context, page, limit int, status string) ([]entity.Coupon, int64, error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 10 }
	return s.couponRepo.GetAll(ctx, page, limit, status)
}

// Approve — admin only.
func (s *CouponService) Approve(ctx context.Context, id uint, adminID uint) error {
	c, err := s.couponRepo.GetByID(ctx, id)
	if err != nil || c == nil {
		return apperror.NewNotFound("coupon not found")
	}
	if c.Status != entity.CouponStatusPending {
		return apperror.NewBadRequest("only pending coupons can be approved")
	}
	return s.couponRepo.UpdateStatus(ctx, id, entity.CouponStatusApproved, adminID)
}

// Reject — admin only.
func (s *CouponService) Reject(ctx context.Context, id uint, adminID uint) error {
	c, err := s.couponRepo.GetByID(ctx, id)
	if err != nil || c == nil {
		return apperror.NewNotFound("coupon not found")
	}
	if c.Status != entity.CouponStatusPending {
		return apperror.NewBadRequest("only pending coupons can be rejected")
	}
	return s.couponRepo.UpdateStatus(ctx, id, entity.CouponStatusRejected, adminID)
}

// ApplyToOrder records coupon usage and queues affiliate payout if applicable.
// Called after payment is confirmed.
func (s *CouponService) ApplyToOrder(ctx context.Context, couponID, orderID, userID uint, orderTotalCents int64) error {
	c, err := s.couponRepo.GetByID(ctx, couponID)
	if err != nil || c == nil {
		return apperror.NewNotFound("coupon not found")
	}

	discount := c.ApplyDiscount(orderTotalCents)
	commission := c.CommissionAmount(orderTotalCents)

	usage := &entity.CouponUsage{
		CouponID:             couponID,
		OrderID:              orderID,
		UserID:               userID,
		DiscountAppliedCents: discount,
		CommissionUSDCents:   commission,
	}
	if err := s.couponUsageRepo.Create(ctx, usage); err != nil {
		return apperror.NewInternal(err, "failed to record coupon usage")
	}
	_ = s.couponRepo.IncrementUsed(ctx, couponID)

	// Queue affiliate payout if coupon has commission and a creator.
	if commission > 0 && c.CreatorID != 0 {
		orderIDCopy := orderID
		payout := &entity.Payout{
			RecipientID: c.CreatorID,
			Type:        entity.PayoutTypeAffiliateCommission,
			AmountCents: commission,
			Status:      entity.PayoutStatusPending,
			OrderID:     &orderIDCopy,
		}
		_ = s.payoutRepo.Create(ctx, payout)
	}
	return nil
}

func validCouponType(t string) bool {
	switch t {
	case entity.CouponTypeDiscount, entity.CouponTypeAffiliate, entity.CouponTypeBoth:
		return true
	}
	return false
}