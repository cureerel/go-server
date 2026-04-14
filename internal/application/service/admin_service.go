// internal/application/service/superadmin_service.go
package service

import (
	"context"

	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/domain/repository"
	"gorm.io/gorm"
)

type SuperAdminService struct {
	userRepo          repository.UserRepository
	db                *gorm.DB
}

func NewSuperAdminService(
	userRepo    repository.UserRepository,
	db          *gorm.DB,
) *SuperAdminService {
	return &SuperAdminService{userRepo: userRepo}
}




//  Platform Stats

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

// ── helpers

var roleRankMap = map[string]int{
	entity.RoleUser:       1,
	entity.RolePartner:    2,
	entity.RoleAdmin: 3,
}

func validRole(r string) bool {
	_, ok := roleRankMap[r]
	return ok
}

