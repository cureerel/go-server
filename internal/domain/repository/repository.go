// internal/domain/repository/repository.go
package repository

import (
	"context"

	"github.com/cureerel/cserver/internal/domain/entity"
)

// ── Ticket ────────────────────────────────────────────────────

type TicketFilter struct {
	Page     int
	Limit    int
	Status   string
	UserID   *uint
	Priority string
}

type TicketRepository interface {
	Create(ctx context.Context, ticket *entity.Ticket) error
	GetByID(ctx context.Context, id uint) (*entity.Ticket, error)
	GetAll(ctx context.Context, filter TicketFilter) ([]entity.Ticket, int64, error)
	Update(ctx context.Context, ticket *entity.Ticket) error
	Close(ctx context.Context, id uint) error
	Assign(ctx context.Context, id uint, workerID uint) error
	CreateMessage(ctx context.Context, msg *entity.TicketMessage) error
	GetMessages(ctx context.Context, ticketID uint) ([]entity.TicketMessage, error)
}

// ── Coupon ────────────────────────────────────────────────────

type CouponRepository interface {
	Create(ctx context.Context, coupon *entity.Coupon) error
	GetByID(ctx context.Context, id uint) (*entity.Coupon, error)
	GetByCode(ctx context.Context, code string) (*entity.Coupon, error)
	GetAll(ctx context.Context, page, limit int, status string) ([]entity.Coupon, int64, error)
	UpdateStatus(ctx context.Context, id uint, status string, approvedBy uint) error
	IncrementUsed(ctx context.Context, id uint) error
}

// ── Payout ────────────────────────────────────────────────────

type PayoutRepository interface {
	Create(ctx context.Context, payout *entity.Payout) error
	GetByID(ctx context.Context, id uint) (*entity.Payout, error)
	GetAll(ctx context.Context, page, limit int, status string) ([]entity.Payout, int64, error)
	GetByRecipient(ctx context.Context, recipientID uint, page, limit int) ([]entity.Payout, int64, error)
	MarkPaid(ctx context.Context, id uint, paidBy uint, reference string) error
}

// ── UpgradeRequest ────────────────────────────────────────────

type UpgradeRequestRepository interface {
	Create(ctx context.Context, req *entity.UpgradeRequest) error
	GetByID(ctx context.Context, id uint) (*entity.UpgradeRequest, error)
	GetPending(ctx context.Context, page, limit int) ([]entity.UpgradeRequest, int64, error)
	GetByUser(ctx context.Context, userID uint) (*entity.UpgradeRequest, error)
	Review(ctx context.Context, id uint, status string, reviewedBy uint) error
}