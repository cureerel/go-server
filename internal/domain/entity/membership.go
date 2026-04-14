package entity

import "time"

type MembershipPlan string
type MembershipStatus string

const (
	PlanFree  MembershipPlan = "free"
	PlanBasic MembershipPlan = "basic"
	PlanPro   MembershipPlan = "pro"
)

const (
	MembershipActive    MembershipStatus = "active"
	MembershipCancelled MembershipStatus = "cancelled"
	MembershipExpired   MembershipStatus = "expired"
)

type Membership struct {
	ID        uint             `json:"id"`
	UserID    uint             `json:"user_id"`
	Plan      MembershipPlan   `json:"plan"`
	Status    MembershipStatus `json:"status"`
	StartsAt  time.Time        `json:"starts_at"`
	ExpiresAt time.Time        `json:"expires_at"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

func (m *Membership) IsActive() bool {
	return m.Status == MembershipActive && time.Now().Before(m.ExpiresAt)
}

func (m *Membership) IsExpired() bool {
	return time.Now().After(m.ExpiresAt)
}
