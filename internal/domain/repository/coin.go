package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrInsufficientCoins = errors.New("insufficient coins")

type CoinLedgerEntry struct {
	ID           int64
	UserID       uint
	Delta        int64
	BalanceAfter int64
	Reason       string
	RefType      string
	RefID        *uint
}

type CoinRepository interface {
	GetBalance(ctx context.Context, userID uint) (int64, error)
	// Credit adds coins (e.g. fiat top-up confirmed). Runs inside optional tx.
	Credit(ctx context.Context, tx *gorm.DB, userID uint, amount int64, reason string, refType string, refID *uint) (balanceAfter int64, err error)
	// Debit subtracts coins; returns ErrInsufficientBalance if needed.
	Debit(ctx context.Context, tx *gorm.DB, userID uint, amount int64, reason string, refType string, refID *uint) (balanceAfter int64, err error)
	HasBlogUnlock(ctx context.Context, userID, blogID uint) (bool, error)
	AddBlogUnlock(ctx context.Context, tx *gorm.DB, userID, blogID uint, coinsSpent int64) error
}
