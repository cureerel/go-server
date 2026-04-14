package models

import "time"

type UserWallet struct {
	UserID    uint      `gorm:"primaryKey"`
	Balance   int64     `gorm:"not null;default:0"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (UserWallet) TableName() string { return "user_wallets" }

type CoinLedger struct {
	ID           int64  `gorm:"primaryKey"`
	UserID       uint   `gorm:"not null;index"`
	Delta        int64  `gorm:"not null"`
	BalanceAfter int64  `gorm:"not null"`
	Reason       string `gorm:"size:50;not null"`
	RefType      string `gorm:"size:50"`
	RefID        *uint
	Meta         []byte `gorm:"type:jsonb"`
	CreatedAt    time.Time
}

func (CoinLedger) TableName() string { return "coin_ledger" }

type BlogUnlock struct {
	UserID     uint      `gorm:"primaryKey"`
	BlogID     uint      `gorm:"primaryKey"`
	CoinsSpent int64     `gorm:"not null"`
	CreatedAt  time.Time `gorm:"not null"`
}

func (BlogUnlock) TableName() string { return "blog_unlocks" }
