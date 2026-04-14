// internal/application/service/admin_service.go
package service

import (
	"context"
	"fmt"

	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/domain/repository"
	"gorm.io/gorm"
)

type AdminService struct {
	userRepo repository.UserRepository
	db       *gorm.DB
}

func NewAdminService(
	userRepo repository.UserRepository,
	db *gorm.DB,
) *AdminService {
	return &AdminService{userRepo: userRepo, db: db}
}

// Platform stats

type PlatformStats struct {
	TotalUsers      int64            `json:"total_users"`
	ByRole          map[string]int64 `json:"by_role"`
	TotalBlogs      int64            `json:"total_blogs"`
	TotalServices   int64            `json:"total_services"`
	TotalOrders     int64            `json:"total_orders"`
	RevenueCents    int64            `json:"revenue_cents"`
	TotalTickets    int64            `json:"total_tickets"`
	OpenTickets     int64            `json:"open_tickets"`
	PendingCoupons  int64            `json:"pending_coupons"`
	PendingUpgrades int64            `json:"pending_upgrades"`
}

func (s *AdminService) PlatformStats(ctx context.Context) (*PlatformStats, error) {
	var stats PlatformStats
	stats.ByRole = make(map[string]int64)

	// Role counts
	type roleCount struct {
		Role  string `gorm:"column:role"`
		Count int64  `gorm:"column:count"`
	}
	var roleCounts []roleCount
	s.db.WithContext(ctx).Raw(`SELECT role, COUNT(*) AS count FROM users WHERE deleted_at IS NULL GROUP BY role`).Scan(&roleCounts)
	for _, rc := range roleCounts {
		stats.TotalUsers += rc.Count
		stats.ByRole[rc.Role] = rc.Count
	}

	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM blogs WHERE deleted_at IS NULL`).Scan(&stats.TotalBlogs)
	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM services WHERE deleted_at IS NULL`).Scan(&stats.TotalServices)
	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM orders`).Scan(&stats.TotalOrders)
	s.db.WithContext(ctx).Raw(`SELECT COALESCE(SUM(total_cents),0) FROM orders WHERE status = 'paid'`).Scan(&stats.RevenueCents)
	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM tickets`).Scan(&stats.TotalTickets)
	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM tickets WHERE status = 'open'`).Scan(&stats.OpenTickets)
	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM coupons WHERE status = 'pending'`).Scan(&stats.PendingCoupons)
	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM users WHERE upgrade_requested_at IS NOT NULL AND role = 'user' AND deleted_at IS NULL`).Scan(&stats.PendingUpgrades)

	return &stats, nil
}

// ChangeRole updates a user's role (admin only).
// Only valid roles (user, partner, admin) are accepted.
func (s *AdminService) ChangeRole(ctx context.Context, userID uint, newRole string) error {
	if !validRole(newRole) {
		return fmt.Errorf("invalid role: must be one of user, partner, admin")
	}
	return s.userRepo.UpdateRole(ctx, userID, newRole)
}

var roleRankMap = map[string]int{
	entity.RoleUser:    1,
	entity.RolePartner: 2,
	entity.RoleAdmin:   3,
}

func validRole(r string) bool {
	_, ok := roleRankMap[r]
	return ok
}
