// internal/domain/entity/user.go
package entity

import (
	"fmt"
	"time"
)

// Role constants — single source of truth for role strings
const (
	RoleUser       = "user"
	RoleWriter     = "writer"
	RoleReviewer   = "reviewer"
	RolePartner    = "partner"
	RoleWorker     = "worker"
	RoleAdmin      = "admin"
	RoleSuperAdmin = "superadmin"
)

type User struct {
	ID                  uint       `json:"id"`
	Name                string     `json:"name"`
	Email               string     `json:"email"`
	Password            string     `json:"-"`
	PasswordHash        string     `json:"-"`
	Role                string     `json:"role"`
	IsActive            bool       `json:"is_active"`
	IsVerified          bool       `json:"is_verified"`
	UpgradeRequestedAt  *time.Time `json:"upgrade_requested_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

func (u *User) GetID() string {
	return fmt.Sprintf("%d", u.ID)
}

func (u *User) MigratePassword() {
	if u.Password != "" && u.PasswordHash == "" {
		u.PasswordHash = u.Password
		u.Password = ""
	}
}

// HasRole returns true if the user's role is at least the given role.
// Order: user < writer < partner < worker < admin < superadmin.
// Reviewer is editorial-only: does not inherit writer/partner routes via this ladder.
func (u *User) HasRole(role string) bool {
	if u.Role == RoleReviewer {
		return role == RoleReviewer
	}
	order := map[string]int{
		RoleUser:       1,
		RoleWriter:     2,
		RolePartner:    3,
		RoleWorker:     4,
		RoleAdmin:      5,
		RoleSuperAdmin: 6,
	}
	// Reviewer is not in the ladder; only explicit RequireAnyRole should grant reviewer APIs.
	if role == RoleReviewer {
		return false
	}
	return order[u.Role] >= order[role]
}

func (u *User) IsSuperAdmin() bool { return u.Role == RoleSuperAdmin }
func (u *User) IsAdmin() bool      { return u.Role == RoleAdmin || u.IsSuperAdmin() }
func (u *User) IsPartner() bool    { return u.Role == RolePartner }
func (u *User) IsWriter() bool     { return u.Role == RoleWriter || u.IsPartner() }
func (u *User) IsWorker() bool     { return u.Role == RoleWorker }