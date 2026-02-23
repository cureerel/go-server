package models

import (
    "time"
    "github.com/cureerel/gotemplate/internal/domain/entity"
)

type Membership struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"not null;uniqueIndex"`
    Plan      string    `gorm:"not null;size:20"`
    Status    string    `gorm:"not null;size:20;default:'active'"`
    StartsAt  time.Time `gorm:"not null"`
    ExpiresAt time.Time `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (Membership) TableName() string {
    return "memberships"
}

func (m *Membership) ToDomain() *entity.Membership {
    return &entity.Membership{
        ID:        m.ID,
        UserID:    m.UserID,
        Plan:      entity.MembershipPlan(m.Plan),
        Status:    entity.MembershipStatus(m.Status),
        StartsAt:  m.StartsAt,
        ExpiresAt: m.ExpiresAt,
        CreatedAt: m.CreatedAt,
        UpdatedAt: m.UpdatedAt,
    }
}

func MembershipFromDomain(e *entity.Membership) *Membership {
    return &Membership{
        ID:        e.ID,
        UserID:    e.UserID,
        Plan:      string(e.Plan),
        Status:    string(e.Status),
        StartsAt:  e.StartsAt,
        ExpiresAt: e.ExpiresAt,
        CreatedAt: e.CreatedAt,
        UpdatedAt: e.UpdatedAt,
    }
}