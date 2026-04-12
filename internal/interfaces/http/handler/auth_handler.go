// internal/interfaces/http/handler/auth_handler.go
package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cureerel/cserver/internal/application/service"
	"github.com/cureerel/cserver/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
	userService *service.UserService
	otpService  *service.OTPService
	otpExpiry   int // minutes, for response hint
}

func NewAuthHandler(
	authService *service.AuthService,
	userService *service.UserService,
	otpService *service.OTPService,
	otpExpiry int,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
		otpService:  otpService,
		otpExpiry:   otpExpiry,
	}
}

func (h *AuthHandler) RegisterInit(c *gin.Context) {
	var req dto.RegisterInitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.otpService.SendRegisterOTP(c.Request.Context(), req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": dto.RegisterInitResponse{
			Message:   "Verification code sent to " + req.Email,
			ExpiresIn: h.otpExpiry,
		},
	})
}

// RegisterVerify godoc
// POST /api/auth/register/verify
// Verifies OTP + creates user + returns tokens (flat response).
func (h *AuthHandler) RegisterVerify(c *gin.Context) {
	var req dto.RegisterVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Step 1: verify OTP
	if err := h.otpService.VerifyRegisterOTP(c.Request.Context(), req.Email, req.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Step 2: create user (OTP already verified, is_verified=true)
	user, err := h.authService.Signup(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Step 3: auto-login — return tokens immediately
	meta := service.LoginMeta{
		UserAgent: c.GetHeader("User-Agent"),
		IPAddress: c.ClientIP(),
	}
	token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password, meta)
	if err != nil {
		// User created but login failed — not critical, tell them to login manually
		c.JSON(http.StatusCreated, gin.H{
			"data": dto.SignupResponse{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
			},
			"message": "Account created. Please log in.",
		})
		return
	}

	// FLAT response: { data: { access_token, refresh_token, expires_at, user } }
	c.JSON(http.StatusCreated, gin.H{
		"data": dto.LoginResponse{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			ExpiresAt:    token.ExpiresAt.Format(time.RFC3339),
			User:         toUserResponse(user),
		},
	})
}

// ── Password Reset ────────────────────────────────────────────

// PasswordResetInit godoc
// POST /api/auth/password/reset/init
func (h *AuthHandler) PasswordResetInit(c *gin.Context) {
	var req dto.PasswordResetInitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Always return success — prevents email enumeration
	_ = h.otpService.SendResetOTP(c.Request.Context(), req.Email)

	c.JSON(http.StatusOK, gin.H{
		"data": dto.PasswordResetInitResponse{
			Message:   "If that email exists, a reset code has been sent.",
			ExpiresIn: h.otpExpiry,
		},
	})
}

// PasswordResetVerify godoc
// POST /api/auth/password/reset/verify
func (h *AuthHandler) PasswordResetVerify(c *gin.Context) {
	var req dto.PasswordResetVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.otpService.VerifyResetOTP(c.Request.Context(), req.Email, req.Code, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully. Please log in."})
}

// ── Change Password (authenticated) ──────────────────────────

// ChangePassword godoc
// POST /api/auth/password/change
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(fmt.Sprintf("%v", userID), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.authService.ChangePassword(c.Request.Context(), uint(id), req.CurrentPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully."})
}

// ── Existing handlers ────────────────────────────────────────

func (h *AuthHandler) Signup(c *gin.Context) {
	var req dto.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.authService.Signup(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"data": dto.SignupResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	meta := service.LoginMeta{
		UserAgent: c.GetHeader("User-Agent"),
		IPAddress: c.ClientIP(),
	}
	token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password, meta)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Fetch user to include in login response
	user, err := h.userService.GetByEmail(c.Request.Context(), req.Email)
	if err != nil || user == nil {
		// Fallback: return tokens without user (frontend will call getMe)
		c.JSON(http.StatusOK, gin.H{
			"data": dto.AuthResponse{
				AccessToken:  token.AccessToken,
				RefreshToken: token.RefreshToken,
				ExpiresAt:    token.ExpiresAt.Format(time.RFC3339),
			},
		})
		return
	}

	// FLAT response: { data: { access_token, refresh_token, expires_at, user } }
	c.JSON(http.StatusOK, gin.H{
		"data": dto.LoginResponse{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			ExpiresAt:    token.ExpiresAt.Format(time.RFC3339),
			User:         toUserResponse(user),
		},
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	meta := service.LoginMeta{
		UserAgent: c.GetHeader("User-Agent"),
		IPAddress: c.ClientIP(),
	}
	token, err := h.authService.Refresh(c.Request.Context(), req.RefreshToken, meta)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": dto.AuthResponse{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			ExpiresAt:    token.ExpiresAt.Format(time.RFC3339),
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accessToken := c.GetHeader("Authorization")
	if len(accessToken) > 7 {
		accessToken = accessToken[7:]
	}
	if err := h.authService.Logout(c.Request.Context(), accessToken, req.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}