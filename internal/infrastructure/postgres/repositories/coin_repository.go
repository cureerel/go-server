package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/cureerel/cserver/internal/domain/repository"
	"github.com/cureerel/cserver/internal/infrastructure/postgres/models"
	"gorm.io/gorm"
)

type coinRepository struct{ db *gorm.DB }

func NewCoinRepository(db *gorm.DB) repository.CoinRepository {
	return &coinRepository{db: db}
}

func dbOrTx(db *gorm.DB, tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return db
}

func (r *coinRepository) GetBalance(ctx context.Context, userID uint) (int64, error) {
	var w models.UserWallet
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&w).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return w.Balance, nil
}

func (r *coinRepository) Credit(ctx context.Context, tx *gorm.DB, userID uint, amount int64, reason, refType string, refID *uint) (int64, error) {
	if amount <= 0 {
		return 0, errors.New("credit amount must be positive")
	}
	conn := dbOrTx(r.db, tx)
	if tx != nil {
		return applyCredit(ctx, tx, userID, amount, reason, refType, refID)
	}
	var after int64
	err := conn.WithContext(ctx).Transaction(func(tx2 *gorm.DB) error {
		var err error
		after, err = applyCredit(ctx, tx2, userID, amount, reason, refType, refID)
		return err
	})
	return after, err
}

func applyCredit(ctx context.Context, tx *gorm.DB, userID uint, amount int64, reason, refType string, refID *uint) (int64, error) {
	if err := tx.WithContext(ctx).Exec(`
		INSERT INTO user_wallets (user_id, balance, updated_at) VALUES (?, ?, NOW())
		ON CONFLICT (user_id) DO NOTHING`, userID, int64(0)).Error; err != nil {
		return 0, err
	}
	if err := tx.WithContext(ctx).Exec(`
		UPDATE user_wallets SET balance = balance + ?, updated_at = NOW() WHERE user_id = ?`, amount, userID).Error; err != nil {
		return 0, err
	}
	var after int64
	if err := tx.WithContext(ctx).Raw(`SELECT balance FROM user_wallets WHERE user_id = ?`, userID).Scan(&after).Error; err != nil {
		return 0, err
	}
	meta, _ := json.Marshal(map[string]any{"ref_type": refType})
	le := models.CoinLedger{
		UserID:       userID,
		Delta:        amount,
		BalanceAfter: after,
		Reason:       reason,
		RefType:      refType,
		RefID:        refID,
		Meta:         meta,
		CreatedAt:    time.Now().UTC(),
	}
	if err := tx.WithContext(ctx).Create(&le).Error; err != nil {
		return 0, err
	}
	return after, nil
}

func (r *coinRepository) Debit(ctx context.Context, tx *gorm.DB, userID uint, amount int64, reason, refType string, refID *uint) (int64, error) {
	if amount <= 0 {
		return 0, errors.New("debit amount must be positive")
	}
	conn := dbOrTx(r.db, tx)
	if tx != nil {
		return applyDebit(ctx, tx, userID, amount, reason, refType, refID)
	}
	var after int64
	err := conn.WithContext(ctx).Transaction(func(tx2 *gorm.DB) error {
		var err error
		after, err = applyDebit(ctx, tx2, userID, amount, reason, refType, refID)
		return err
	})
	return after, err
}

func applyDebit(ctx context.Context, tx *gorm.DB, userID uint, amount int64, reason, refType string, refID *uint) (int64, error) {
	res := tx.WithContext(ctx).Exec(`
		UPDATE user_wallets SET balance = balance - ?, updated_at = NOW()
		WHERE user_id = ? AND balance >= ?`, amount, userID, amount)
	if res.Error != nil {
		return 0, res.Error
	}
	if res.RowsAffected == 0 {
		return 0, repository.ErrInsufficientCoins
	}
	var after int64
	if err := tx.WithContext(ctx).Raw(`SELECT balance FROM user_wallets WHERE user_id = ?`, userID).Scan(&after).Error; err != nil {
		return 0, err
	}
	meta, _ := json.Marshal(map[string]any{"ref_type": refType})
	le := models.CoinLedger{
		UserID:       userID,
		Delta:        -amount,
		BalanceAfter: after,
		Reason:       reason,
		RefType:      refType,
		RefID:        refID,
		Meta:         meta,
		CreatedAt:    time.Now().UTC(),
	}
	if err := tx.WithContext(ctx).Create(&le).Error; err != nil {
		return 0, err
	}
	return after, nil
}

func (r *coinRepository) HasBlogUnlock(ctx context.Context, userID, blogID uint) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.BlogUnlock{}).
		Where("user_id = ? AND blog_id = ?", userID, blogID).Count(&n).Error
	return n > 0, err
}

func (r *coinRepository) AddBlogUnlock(ctx context.Context, tx *gorm.DB, userID, blogID uint, coinsSpent int64) error {
	conn := dbOrTx(r.db, tx)
	u := models.BlogUnlock{
		UserID:     userID,
		BlogID:     blogID,
		CoinsSpent: coinsSpent,
		CreatedAt:  time.Now().UTC(),
	}
	return conn.WithContext(ctx).Create(&u).Error
}
