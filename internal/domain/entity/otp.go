// internal/domain/entity/otp.go
package entity

import "time"

const (
	OTPTypeRegister = "register"
	OTPTypeReset    = "reset"
)

type OTP struct {
	ID        uint      `json:"id"`
	UserID    *uint     `json:"user_id,omitempty"` 
	Email     string    `json:"email"`
	Code      string    `json:"-"`
	Type      string    `json:"type"`
	Used      bool      `json:"used"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (o *OTP) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

func (o *OTP) IsValid() bool {
	return !o.Used && !o.IsExpired()
}