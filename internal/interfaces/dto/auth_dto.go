// internal/interfaces/dto/auth_dto.go
package dto

// ── Registration (2-step OTP flow) ───────────────────────────

type RegisterInitRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type RegisterInitResponse struct {
	Message   string `json:"message"`
	ExpiresIn int    `json:"expires_in_minutes"`
}

type RegisterVerifyRequest struct {
	Name     string `json:"name"     binding:"required,min=2,max=100"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Code     string `json:"code"     binding:"required,len=6"`
}



type PasswordResetInitRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type PasswordResetInitResponse struct {
	Message   string `json:"message"`
	ExpiresIn int    `json:"expires_in_minutes"`
}

type PasswordResetVerifyRequest struct {
	Email       string `json:"email"        binding:"required,email"`
	Code        string `json:"code"         binding:"required,len=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ── Login ─────────────────────────────────────────────────────

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// ── Token responses ───────────────────────────────────────────

// AuthResponse is used for token-only responses (refresh).
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}

// LoginResponse is the flat response for login & register/verify.
// Frontend expects { access_token, refresh_token, user } at the top level.
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresAt    string       `json:"expires_at"`
	User         UserResponse `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ── Change password (authenticated) ──────────────────────────

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password"     binding:"required,min=6"`
}

// ── Signup (legacy direct) ───────────────────────────────────

type SignupRequest struct {
	Name     string `json:"name"     binding:"required,min=2,max=100"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type SignupResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}