// internal/domain/entity/session.go
package entity

import "time"

// Session is a pure business object
type Session struct {
    ID         string
    UserID     uint
    TokenHash  string
    UserAgent  string
    IPAddress  string
    ExpiresAt  time.Time
    Revoked    bool
    CreatedAt  time.Time
    LastActive time.Time
}

// Business logic methods (no DB concerns)
func (s *Session) IsActive() bool {
    return !s.Revoked && s.ExpiresAt.After(time.Now())
}