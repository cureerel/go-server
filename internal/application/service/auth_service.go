// internal/application/service/auth_service.go
// CHANGES from original:
//   - Signup now sets is_verified=true only after OTP is confirmed.
//     The OTP verification flow calls userRepo.UpdateVerified().
//   - hashPassword helper extracted here so otp_service can reuse it.
//   - Everything else unchanged.
package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/domain/repository"
	"github.com/cureerel/cserver/pkg/apperror"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo      repository.UserRepository
	authRepo      repository.AuthRepository
	sessionRepo   repository.SessionRepository
	accessSecret  []byte
	refreshSecret []byte
	tokenHashKey  []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	TokenHashKey  string
}

func NewAuthService(
	userRepo repository.UserRepository,
	authRepo repository.AuthRepository,
	sessionRepo repository.SessionRepository,
	cfg JWTConfig,
) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		authRepo:      authRepo,
		sessionRepo:   sessionRepo,
		accessSecret:  []byte(cfg.AccessSecret),
		refreshSecret: []byte(cfg.RefreshSecret),
		tokenHashKey:  []byte(cfg.TokenHashKey),
		accessTTL:     15 * time.Minute,
		refreshTTL:    7 * 24 * time.Hour,
	}
}

type LoginMeta struct {
	UserAgent string
	IPAddress string
}

// hashPassword is a shared helper used by auth and OTP services.
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Signup creates the user record after OTP has already been verified.
// The handler flow is:
//  1. POST /auth/register/init   → OTPService.SendRegisterOTP
//  2. POST /auth/register/verify → OTPService.VerifyRegisterOTP → AuthService.Signup
func (s *AuthService) Signup(ctx context.Context, name, email, password string) (*entity.User, error) {
	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to check existing user")
	}
	if existing != nil {
		return nil, apperror.NewBadRequest("email already registered")
	}

	hash, err := hashPassword(password)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to hash password")
	}

	user := &entity.User{
		Name:         name,
		Email:        email,
		PasswordHash: hash,
		Role:         entity.RoleUser,
		IsActive:     true,
		IsVerified:   true, // OTP already verified before this is called
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, apperror.NewInternal(err, "failed to create user")
	}

	return user, nil
}

// Login authenticates and returns token pair.
// Unverified users are blocked.
func (s *AuthService) Login(ctx context.Context, email, password string, meta LoginMeta) (*entity.AuthToken, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, apperror.NewBadRequest("invalid credentials")
	}

	if !user.IsVerified {
		return nil, apperror.NewBadRequest("email not verified — check your inbox for the verification code")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, apperror.NewBadRequest("invalid credentials")
	}

	return s.generateTokenPair(ctx, user.ID, user.Email, user.Role, meta)
}

func (s *AuthService) generateTokenPair(
	ctx context.Context,
	userID uint,
	email, role string,
	meta LoginMeta,
) (*entity.AuthToken, error) {
	now := time.Now()

	accessClaims := jwt.MapClaims{
		"user_id":    strconv.Itoa(int(userID)),
		"email":      email,
		"role":       role,
		"token_type": "access",
		"exp":        now.Add(s.accessTTL).Unix(),
		"iat":        now.Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessToken.SignedString(s.accessSecret)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to sign access token")
	}

	jti := uuid.New().String()
	refreshClaims := jwt.MapClaims{
		"user_id":    strconv.Itoa(int(userID)),
		"token_type": "refresh",
		"exp":        now.Add(s.refreshTTL).Unix(),
		"iat":        now.Unix(),
		"jti":        jti,
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := refreshToken.SignedString(s.refreshSecret)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to sign refresh token")
	}

	tokenHash := hashToken(refreshStr, s.tokenHashKey)

	refreshEntity := &entity.RefreshToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: now.Add(s.refreshTTL),
	}
	if err := s.authRepo.SaveRefreshToken(ctx, refreshEntity); err != nil {
		return nil, apperror.NewInternal(err, "failed to save refresh token")
	}

	session := &entity.Session{
		UserID:     userID,
		TokenHash:  tokenHash,
		UserAgent:  meta.UserAgent,
		IPAddress:  meta.IPAddress,
		ExpiresAt:  now.Add(s.refreshTTL),
		LastActive: now,
	}
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, apperror.NewInternal(err, "failed to create session")
	}

	return &entity.AuthToken{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
		ExpiresAt:    now.Add(s.accessTTL),
	}, nil
}

func hashToken(token string, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

// Refresh, Logout, LogoutAll, GetActiveSessions, ValidateAccessToken
// are unchanged from the original — keeping them here for completeness.

func (s *AuthService) Refresh(ctx context.Context, refreshToken string, meta LoginMeta) (*entity.AuthToken, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.refreshSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, apperror.NewBadRequest("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["token_type"] != "refresh" {
		return nil, apperror.NewBadRequest("invalid token type")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, apperror.NewBadRequest("invalid user id in token")
	}

	tokenHash := hashToken(refreshToken, s.tokenHashKey)

	stored, err := s.authRepo.GetRefreshToken(ctx, tokenHash)
	if err != nil || stored == nil || stored.Revoked || stored.ExpiresAt.Before(time.Now()) {
		return nil, apperror.NewBadRequest("refresh token revoked or expired")
	}
	if strconv.Itoa(int(stored.UserID)) != userIDStr {
		return nil, apperror.NewBadRequest("token user mismatch")
	}

	session, err := s.sessionRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to look up session")
	}
	if session == nil || !session.IsActive() {
		return nil, apperror.NewBadRequest("session expired or revoked")
	}

	if err := s.authRepo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		return nil, apperror.NewInternal(err, "failed to revoke old refresh token")
	}
	if err := s.sessionRepo.Revoke(ctx, session.ID); err != nil {
		return nil, apperror.NewInternal(err, "failed to revoke old session")
	}

	user, err := s.userRepo.GetByID(ctx, stored.UserID)
	if err != nil || user == nil {
		return nil, apperror.NewBadRequest("user not found")
	}

	return s.generateTokenPair(ctx, user.ID, user.Email, user.Role, meta)
}

func (s *AuthService) Logout(ctx context.Context, accessToken, refreshToken string) error {
	if refreshToken != "" {
		tokenHash := hashToken(refreshToken, s.tokenHashKey)
		if err := s.authRepo.RevokeRefreshToken(ctx, tokenHash); err != nil {
			return apperror.NewInternal(err, "failed to revoke refresh token")
		}
		session, err := s.sessionRepo.GetByTokenHash(ctx, tokenHash)
		if err != nil {
			return apperror.NewInternal(err, "failed to look up session")
		}
		if session != nil {
			if err := s.sessionRepo.Revoke(ctx, session.ID); err != nil {
				return apperror.NewInternal(err, "failed to revoke session")
			}
		}
	}
	if accessToken != "" {
		tokenHash := hashToken(accessToken, s.tokenHashKey)
		if err := s.authRepo.BlacklistToken(ctx, tokenHash, time.Now().Add(s.accessTTL)); err != nil {
			return apperror.NewInternal(err, "failed to blacklist access token")
		}
	}
	return nil
}

func (s *AuthService) LogoutAll(ctx context.Context, userID uint, currentAccessToken string) error {
	if err := s.sessionRepo.RevokeAllByUserID(ctx, userID); err != nil {
		return apperror.NewInternal(err, "failed to revoke all sessions")
	}
	if currentAccessToken != "" {
		tokenHash := hashToken(currentAccessToken, s.tokenHashKey)
		if err := s.authRepo.BlacklistToken(ctx, tokenHash, time.Now().Add(s.accessTTL)); err != nil {
			return apperror.NewInternal(err, "failed to blacklist access token")
		}
	}
	return nil
}

func (s *AuthService) GetActiveSessions(ctx context.Context, userID uint) ([]*entity.Session, error) {
	sessions, err := s.sessionRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch sessions")
	}
	return sessions, nil
}

func (s *AuthService) ValidateAccessToken(tokenString string) (*entity.Claims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.accessSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid access token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid || claims["token_type"] != "access" {
		return nil, fmt.Errorf("invalid access token")
	}

	userIDStr, _ := claims["user_id"].(string)
	email, _ := claims["email"].(string)
	role, _ := claims["role"].(string)

	return &entity.Claims{
		UserID: userIDStr,
		Email:  email,
		Role:   role,
	}, nil
}

// ChangePassword verifies the current password and updates to a new one.
func (s *AuthService) ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return apperror.NewBadRequest("user not found")
	}

	// Check which field has the hash
	hash := user.PasswordHash
	if hash == "" {
		hash = user.Password
	}
	if hash == "" {
		return apperror.NewInternal(nil, "no password hash stored")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(currentPassword)); err != nil {
		return apperror.NewBadRequest("current password is incorrect")
	}

	newHash, err := hashPassword(newPassword)
	if err != nil {
		return apperror.NewInternal(err, "failed to hash new password")
	}

	user.PasswordHash = newHash
	if err := s.userRepo.Update(ctx, user); err != nil {
		return apperror.NewInternal(err, "failed to update password")
	}

	return nil
}