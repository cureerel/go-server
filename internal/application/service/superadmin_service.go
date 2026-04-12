// internal/application/service/superadmin_service.go
package service

import (
	"context"

	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/domain/repository"
	"github.com/cureerel/cserver/pkg/apperror"
	"gorm.io/gorm"
)

type SuperAdminService struct {
	userRepo          repository.UserRepository
	upgradeRepo       repository.UpgradeRequestRepository
	db                *gorm.DB
}

func NewSuperAdminService(
	userRepo    repository.UserRepository,
	upgradeRepo repository.UpgradeRequestRepository,
	db          *gorm.DB,
) *SuperAdminService {
	return &SuperAdminService{userRepo: userRepo, upgradeRepo: upgradeRepo, db: db}
}

// ── Role Management ───────────────────────────────────────────

// SetRole directly sets any user's role. Superadmin only.
func (s *SuperAdminService) SetRole(ctx context.Context, userID uint, role string, callerID uint) error {
	if !validRole(role) {
		return apperror.NewBadRequest("invalid role: " + role)
	}
	if userID == callerID {
		return apperror.NewBadRequest("cannot change your own role")
	}
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return apperror.NewNotFound("user not found")
	}
	// Prevent demoting another superadmin
	if user.Role == entity.RoleSuperAdmin {
		return apperror.NewForbidden("cannot modify another superadmin")
	}
	return s.userRepo.UpdateRole(ctx, userID, role)
}

// ── Upgrade Requests ──────────────────────────────────────────

// RequestUpgrade — any user submits a role upgrade request.
// Only one pending request allowed at a time.
func (s *SuperAdminService) RequestUpgrade(ctx context.Context, userID uint, toRole string) (*entity.UpgradeRequest, error) {
	if !upgradableRole(toRole) {
		return nil, apperror.NewBadRequest("invalid target role")
	}
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, apperror.NewNotFound("user not found")
	}
	if roleRankMap[user.Role] >= roleRankMap[toRole] {
		return nil, apperror.NewBadRequest("you already have this role or higher")
	}
	existing, err := s.upgradeRepo.GetByUser(ctx, userID)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to check existing request")
	}
	if existing != nil && existing.Status == entity.UpgradeStatusPending {
		return nil, apperror.NewConflict("you already have a pending upgrade request")
	}
	req := &entity.UpgradeRequest{
		UserID:   userID,
		FromRole: user.Role,
		ToRole:   toRole,
		Status:   entity.UpgradeStatusPending,
	}
	if err := s.upgradeRepo.Create(ctx, req); err != nil {
		return nil, apperror.NewInternal(err, "failed to create upgrade request")
	}
	return req, nil
}

// GetMyUpgradeRequest returns the caller's own request.
func (s *SuperAdminService) GetMyUpgradeRequest(ctx context.Context, userID uint) (*entity.UpgradeRequest, error) {
	req, err := s.upgradeRepo.GetByUser(ctx, userID)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch upgrade request")
	}
	if req == nil {
		return nil, apperror.NewNotFound("no upgrade request found")
	}
	return req, nil
}

// GetPendingUpgrades — superadmin lists pending requests.
func (s *SuperAdminService) GetPendingUpgrades(ctx context.Context, page, limit int) ([]entity.UpgradeRequest, int64, error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 10 }
	return s.upgradeRepo.GetPending(ctx, page, limit)
}

// ReviewUpgrade — superadmin approves or rejects.
// On approval, the user's role is updated atomically.
func (s *SuperAdminService) ReviewUpgrade(ctx context.Context, reqID uint, approve bool, reviewerID uint) error {
	req, err := s.upgradeRepo.GetByID(ctx, reqID)
	if err != nil || req == nil {
		return apperror.NewNotFound("upgrade request not found")
	}
	if req.Status != entity.UpgradeStatusPending {
		return apperror.NewBadRequest("request is already reviewed")
	}
	status := entity.UpgradeStatusRejected
	if approve {
		status = entity.UpgradeStatusApproved
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.upgradeRepo.Review(ctx, reqID, status, reviewerID); err != nil {
			return err
		}
		if approve {
			if err := s.userRepo.UpdateRole(ctx, req.UserID, req.ToRole); err != nil {
				return err
			}
		}
		return nil
	})
}

// ── Platform Stats ────────────────────────────────────────────

type PlatformStats struct {
	// Users
	TotalUsers      int64 `json:"total_users"`
	ByRole          map[string]int64 `json:"by_role"`
	// Content
	TotalBlogs      int64 `json:"total_blogs"`
	TotalServices   int64 `json:"total_services"`
	// Commerce
	TotalOrders     int64 `json:"total_orders"`
	RevenueCents    int64 `json:"revenue_cents"`
	TotalPayouts    int64 `json:"total_payouts"`
	PendingPayoutsCents int64 `json:"pending_payouts_cents"`
	// Support
	TotalTickets    int64 `json:"total_tickets"`
	OpenTickets     int64 `json:"open_tickets"`
	// Pending actions
	PendingServices  int64 `json:"pending_services"`
	PendingCoupons   int64 `json:"pending_coupons"`
	PendingUpgrades  int64 `json:"pending_upgrades"`
}

func (s *SuperAdminService) PlatformStats(ctx context.Context) (*PlatformStats, error) {
	var stats PlatformStats
	stats.ByRole = make(map[string]int64)

	// User counts by role
	type roleCount struct {
		Role  string `gorm:"column:role"`
		Count int64  `gorm:"column:count"`
	}
	var roleCounts []roleCount
	s.db.WithContext(ctx).Raw(`SELECT role, COUNT(*) AS count FROM users GROUP BY role`).Scan(&roleCounts)
	for _, rc := range roleCounts {
		stats.TotalUsers += rc.Count
		stats.ByRole[rc.Role] = rc.Count
	}

	// Content
	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM blogs`).Scan(&stats.TotalBlogs)
	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM services`).Scan(&stats.TotalServices)

	// Commerce
	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM orders`).Scan(&stats.TotalOrders)
	s.db.WithContext(ctx).Raw(`SELECT COALESCE(SUM(total_cents),0) FROM orders WHERE status = 'confirmed'`).Scan(&stats.RevenueCents)
	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM payouts`).Scan(&stats.TotalPayouts)
	s.db.WithContext(ctx).Raw(`SELECT COALESCE(SUM(amount_cents),0) FROM payouts WHERE status = 'pending'`).Scan(&stats.PendingPayoutsCents)

	// Support
	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM tickets`).Scan(&stats.TotalTickets)
	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM tickets WHERE status = 'open'`).Scan(&stats.OpenTickets)

	// Pending actions
	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM services WHERE status = 'pending'`).Scan(&stats.PendingServices)
	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM coupons WHERE status = 'pending'`).Scan(&stats.PendingCoupons)
	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM upgrade_requests WHERE status = 'pending'`).Scan(&stats.PendingUpgrades)

	return &stats, nil
}

// ── helpers ───────────────────────────────────────────────────

var roleRankMap = map[string]int{
	entity.RoleUser:       1,
	entity.RoleWriter:     2,
	entity.RoleReviewer:   3,
	entity.RolePartner:    4,
	entity.RoleWorker:     5,
	entity.RoleAdmin:      6,
	entity.RoleSuperAdmin: 7,
}

func validRole(r string) bool {
	_, ok := roleRankMap[r]
	return ok
}

func upgradableRole(r string) bool {
	switch r {
	case entity.RoleWriter, entity.RoleReviewer, entity.RolePartner, entity.RoleWorker, entity.RoleAdmin:
		return true
	}
	return false
}