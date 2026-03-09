// internal/domain/entity/service.go
package entity

import "time"

const (
	ServiceStatusPending  = "pending"
	ServiceStatusApproved = "approved"
	ServiceStatusRejected = "rejected"
	ServiceStatusLive     = "live"
	ServiceStatusPaused   = "paused"
)

type Service struct {
	ID            uint      `json:"id"`
	OwnerID       uint      `json:"owner_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	PriceUSDCents int64     `json:"price_usd_cents"`
	Status        string    `json:"status"`
	CoverImageURL string    `json:"cover_image_url,omitempty"`
	CoverImageKey string    `json:"-"`
	ViewsTotal    int64     `json:"views_total"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (s *Service) PriceUSD() float64 { return float64(s.PriceUSDCents) / 100 }
func (s *Service) IsLive() bool      { return s.Status == ServiceStatusLive }
func (s *Service) IsPending() bool   { return s.Status == ServiceStatusPending }
func (s *Service) IsApproved() bool  { return s.Status == ServiceStatusApproved }

type ServiceView struct {
	ID          uint      `json:"id"`
	ServiceID   uint      `json:"service_id"`
	VisitorHash string    `json:"-"`
	ViewedDate  time.Time `json:"viewed_date"`
}