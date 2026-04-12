// internal/application/service/dashboard_service.go
package service

import (
	"context"

	"github.com/cureerel/cserver/pkg/apperror"
	"gorm.io/gorm"
)

// DashboardService runs raw aggregate queries directly on the DB.
// It intentionally bypasses the repository layer — dashboards are
// read-only projections and don't belong to any single domain.
type DashboardService struct {
	db *gorm.DB
}

func NewDashboardService(db *gorm.DB) *DashboardService {
	return &DashboardService{db: db}
}

// ── User dashboard ────────────────────────────────────────────

type UserDashboard struct {
	TotalOrders      int64 `json:"total_orders"`
	PendingOrders    int64 `json:"pending_orders"`
	CompletedOrders  int64 `json:"completed_orders"`
	TotalSpentCents  int64 `json:"total_spent_cents"`
	OpenTickets      int64 `json:"open_tickets"`
	ResolvedTickets  int64 `json:"resolved_tickets"`
}

func (s *DashboardService) UserDashboard(ctx context.Context, userID uint) (*UserDashboard, error) {
	var d UserDashboard
	err := s.db.WithContext(ctx).Raw(`
		SELECT
			COUNT(*)                                          AS total_orders,
			COUNT(*) FILTER (WHERE status = 'pending')       AS pending_orders,
			COUNT(*) FILTER (WHERE status = 'confirmed')     AS completed_orders,
			COALESCE(SUM(total_cents) FILTER (WHERE status = 'confirmed'), 0) AS total_spent_cents
		FROM orders WHERE user_id = ?`, userID).Scan(&d).Error
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch user dashboard")
	}
	err = s.db.WithContext(ctx).Raw(`
		SELECT
			COUNT(*) FILTER (WHERE status = 'open')     AS open_tickets,
			COUNT(*) FILTER (WHERE status = 'resolved') AS resolved_tickets
		FROM tickets WHERE user_id = ?`, userID).Scan(&d).Error
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch ticket stats")
	}
	return &d, nil
}

// ── Partner dashboard ─────────────────────────────────────────

type PartnerDashboard struct {
	TotalServices    int64 `json:"total_services"`
	LiveServices     int64 `json:"live_services"`
	TotalOrders      int64 `json:"total_orders"`
	TotalRevenueCents int64 `json:"total_revenue_cents"`
	PendingPayouts   int64 `json:"pending_payouts"`
	PaidPayouts      int64 `json:"paid_payouts"`
	TotalCoupons     int64 `json:"total_coupons"`
	ActiveCoupons    int64 `json:"active_coupons"`
}

func (s *DashboardService) PartnerDashboard(ctx context.Context, partnerID uint) (*PartnerDashboard, error) {
	var d PartnerDashboard
	err := s.db.WithContext(ctx).Raw(`
		SELECT
			COUNT(*)                                        AS total_services,
			COUNT(*) FILTER (WHERE status = 'live')        AS live_services
		FROM services WHERE creator_id = ?`, partnerID).Scan(&d).Error
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch service stats")
	}
	err = s.db.WithContext(ctx).Raw(`
		SELECT
			COUNT(o.id)                                        AS total_orders,
			COALESCE(SUM(o.total_cents) FILTER (WHERE o.status = 'confirmed'), 0) AS total_revenue_cents
		FROM orders o
		JOIN services sv ON sv.id = o.service_id
		WHERE sv.creator_id = ?`, partnerID).Scan(&d).Error
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch revenue stats")
	}
	err = s.db.WithContext(ctx).Raw(`
		SELECT
			COUNT(*) FILTER (WHERE status = 'pending') AS pending_payouts,
			COUNT(*) FILTER (WHERE status = 'paid')    AS paid_payouts
		FROM payouts WHERE recipient_id = ?`, partnerID).Scan(&d).Error
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch payout stats")
	}
	err = s.db.WithContext(ctx).Raw(`
		SELECT
			COUNT(*)                                          AS total_coupons,
			COUNT(*) FILTER (WHERE status = 'approved')      AS active_coupons
		FROM coupons WHERE creator_id = ?`, partnerID).Scan(&d).Error
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch coupon stats")
	}
	return &d, nil
}

// ── Writer dashboard ──────────────────────────────────────────

type WriterDashboard struct {
	TotalPosts     int64 `json:"total_posts"`
	PublishedPosts int64 `json:"published_posts"`
	DraftPosts     int64 `json:"draft_posts"`
	TotalViews     int64 `json:"total_views"`
}

func (s *DashboardService) WriterDashboard(ctx context.Context, writerID uint) (*WriterDashboard, error) {
	var d WriterDashboard
	err := s.db.WithContext(ctx).Raw(`
		SELECT
			COUNT(*)                                          AS total_posts,
			COUNT(*) FILTER (WHERE published = true)         AS published_posts,
			COUNT(*) FILTER (WHERE published = false)        AS draft_posts,
			COALESCE(SUM(views_total), 0)                    AS total_views
		FROM blogs WHERE author_id = ?`, writerID).Scan(&d).Error
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch writer dashboard")
	}
	return &d, nil
}

// ── Worker dashboard ──────────────────────────────────────────

type WorkerDashboard struct {
	TotalTickets      int64 `json:"total_tickets"`
	OpenTickets       int64 `json:"open_tickets"`
	InProgressTickets int64 `json:"in_progress_tickets"`
	ResolvedTickets   int64 `json:"resolved_tickets"`
	AssignedToMe      int64 `json:"assigned_to_me"`
}

func (s *DashboardService) WorkerDashboard(ctx context.Context, workerID uint) (*WorkerDashboard, error) {
	var d WorkerDashboard
	err := s.db.WithContext(ctx).Raw(`
		SELECT
			COUNT(*)                                                    AS total_tickets,
			COUNT(*) FILTER (WHERE status = 'open')                    AS open_tickets,
			COUNT(*) FILTER (WHERE status = 'in_progress')             AS in_progress_tickets,
			COUNT(*) FILTER (WHERE status = 'resolved')                AS resolved_tickets,
			COUNT(*) FILTER (WHERE assigned_to = ?)                    AS assigned_to_me
		FROM tickets`, workerID).Scan(&d).Error
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch worker dashboard")
	}
	return &d, nil
}

// ── Admin dashboard ───────────────────────────────────────────

type AdminDashboard struct {
	TotalUsers       int64 `json:"total_users"`
	TotalPartners    int64 `json:"total_partners"`
	TotalOrders      int64 `json:"total_orders"`
	RevenueCents     int64 `json:"revenue_cents"`
	PendingPayouts   int64 `json:"pending_payouts"`
	PayoutAmountCents int64 `json:"payout_amount_cents"`
	OpenTickets      int64 `json:"open_tickets"`
	PendingServices  int64 `json:"pending_services"`
	PendingCoupons   int64 `json:"pending_coupons"`
	PendingUpgrades  int64 `json:"pending_upgrades"`
}

func (s *DashboardService) AdminDashboard(ctx context.Context) (*AdminDashboard, error) {
	var d AdminDashboard
	err := s.db.WithContext(ctx).Raw(`
		SELECT
			COUNT(*)                                          AS total_users,
			COUNT(*) FILTER (WHERE role = 'partner')         AS total_partners
		FROM users`).Scan(&d).Error
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch user stats")
	}
	err = s.db.WithContext(ctx).Raw(`
		SELECT
			COUNT(*)                                                   AS total_orders,
			COALESCE(SUM(total_cents) FILTER (WHERE status = 'confirmed'), 0) AS revenue_cents
		FROM orders`).Scan(&d).Error
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch order stats")
	}
	err = s.db.WithContext(ctx).Raw(`
		SELECT
			COUNT(*) FILTER (WHERE status = 'pending')            AS pending_payouts,
			COALESCE(SUM(amount_cents) FILTER (WHERE status = 'pending'), 0) AS payout_amount_cents
		FROM payouts`).Scan(&d).Error
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch payout stats")
	}
	err = s.db.WithContext(ctx).Raw(`
		SELECT
			COUNT(*) FILTER (WHERE status = 'open')    AS open_tickets,
			COUNT(*) FILTER (WHERE status = 'pending') AS pending_services
		FROM tickets, (SELECT 1) _`).Scan(&d).Error
	// simpler: run separately
	var misc struct {
		PendingServices int64 `gorm:"column:pending_services"`
		PendingCoupons  int64 `gorm:"column:pending_coupons"`
		PendingUpgrades int64 `gorm:"column:pending_upgrades"`
		OpenTickets     int64 `gorm:"column:open_tickets"`
	}
	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) AS pending_services FROM services WHERE status = 'pending'`).Scan(&misc)
	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) AS pending_coupons  FROM coupons  WHERE status = 'pending'`).Scan(&misc)
	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) AS pending_upgrades FROM upgrade_requests WHERE status = 'pending'`).Scan(&misc)
	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) AS open_tickets FROM tickets WHERE status = 'open'`).Scan(&misc)
	d.PendingServices = misc.PendingServices
	d.PendingCoupons  = misc.PendingCoupons
	d.PendingUpgrades = misc.PendingUpgrades
	d.OpenTickets     = misc.OpenTickets
	return &d, nil
}