package service

import (
	"context"
	"errors"

	"github.com/cureerel/cserver/internal/domain/repository"
	"github.com/cureerel/cserver/pkg/apperror"
	"gorm.io/gorm"
)

type CoinService struct {
	db       *gorm.DB
	coinRepo repository.CoinRepository
}

func NewCoinService(db *gorm.DB, coinRepo repository.CoinRepository) *CoinService {
	return &CoinService{db: db, coinRepo: coinRepo}
}

func (s *CoinService) Balance(ctx context.Context, userID uint) (int64, error) {
	return s.coinRepo.GetBalance(ctx, userID)
}

// CreditTopUp records coins purchased via external payment (Razorpay / Stripe / DOD).
func (s *CoinService) CreditTopUp(ctx context.Context, userID uint, coins int64, refType string, refID *uint) error {
	if coins <= 0 {
		return apperror.NewBadRequest("invalid coin amount")
	}
	_, err := s.coinRepo.Credit(ctx, nil, userID, coins, "fiat_topup", refType, refID)
	if err != nil {
		return apperror.NewInternal(err, "failed to credit coins")
	}
	return nil
}

// UnlockBlog spends coins and records access to a paid post.
func (s *CoinService) UnlockBlog(ctx context.Context, userID, blogID uint, price int64) error {
	if price <= 0 {
		return apperror.NewBadRequest("invalid price")
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		_, err := s.coinRepo.Debit(ctx, tx, userID, price, "blog_unlock", "blog", &blogID)
		if err != nil {
			if errors.Is(err, repository.ErrInsufficientCoins) {
				return apperror.NewBadRequest("insufficient coins")
			}
			return err
		}
		return s.coinRepo.AddBlogUnlock(ctx, tx, userID, blogID, price)
	})
}
