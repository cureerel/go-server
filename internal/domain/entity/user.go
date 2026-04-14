// internal/domain/entity/user.go
package entity

import (
	"fmt"
	"time"
)

// Role constants  ( 3 roles only )
const (
	RoleUser    = "user"
	RolePartner = "partner"
	RoleAdmin   = "admin"
)

type User struct {
	ID                 uint       `json:"id"`
	Username           string     `json:"username"`
	Email              string     `json:"email"`
	Password           string     `json:"-"`
	PasswordHash       string     `json:"-"`
	Role               string     `json:"role"`
	FirstName          string     `json:"first_name,omitempty"`
	LastName           string     `json:"last_name,omitempty"`
	Country            string     `json:"country,omitempty"`
	PhoneNumber        string     `json:"phone_number,omitempty"`
	Address            string     `json:"address,omitempty"`
	IsActive           bool       `json:"is_active"`
	IsVerified         bool       `json:"is_verified"`
	UpgradeRequestedAt *time.Time `json:"upgrade_requested_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
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

func (u *User) HasRole(role string) bool {
	ranks := map[string]int{
		RoleUser:    1,
		RolePartner: 2,
		RoleAdmin:   3,
	}
	return ranks[u.Role] >= ranks[role]
}

func (u *User) IsAdmin() bool   { return u.Role == RoleAdmin }
func (u *User) IsPartner() bool { return u.Role == RolePartner }
