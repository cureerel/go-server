package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cureerel/gotemplate/internal/domain/entity"
	"github.com/cureerel/gotemplate/internal/domain/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo      repository.UserRepository
	authRepo      repository.AuthRepository
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
}

func NewAuthService(userRepo repository.UserRepository, authRepo repository.AuthRepository, cfg JWTConfig) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		authRepo:      authRepo,
		accessSecret:  []byte(cfg.AccessSecret),
		refreshSecret: []byte(cfg.RefreshSecret),
		accessTTL:     15 * time.Minute,
		refreshTTL:    7 * 24 * time.Hour,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*entity.AuthToken, error) {
	user, err := s.userRepo.GetByEmailContext(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	hash := user.PasswordHash
	if hash == "" {
		hash = user.Password
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	role := user.Role
	if role == "" {
		role = "user"
	}

	return s.generateTokenPair(user.GetID(), user.Email, role)
}

func (s *AuthService) generateTokenPair(userID, email, role string) (*entity.AuthToken, error) {
	now := time.Now()

	accessClaims := jwt.MapClaims{
		"user_id":    userID,
		"email":      email,
		"role":       role,
		"token_type": "access",
		"exp":        now.Add(s.accessTTL).Unix(),
		"iat":        now.Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString(s.accessSecret)
	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.MapClaims{
		"user_id":    userID,
		"token_type": "refresh",
		"exp":        now.Add(s.refreshTTL).Unix(),
		"iat":        now.Unix(),
		"jti":        uuid.New().String(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString(s.refreshSecret)
	if err != nil {
		return nil, err
	}

	if err := s.authRepo.StoreRefreshToken(context.Background(), userID, refreshString, now.Add(s.refreshTTL)); err != nil {
		return nil, err
	}

	return &entity.AuthToken{
		AccessToken:  accessString,
		RefreshToken: refreshString,
		ExpiresAt:    now.Add(s.accessTTL),
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*entity.AuthToken, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.refreshSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["token_type"] != "refresh" {
		return nil, fmt.Errorf("invalid token type")
	}

	userIDStr, _ := claims["user_id"].(string)

	storedUserID, err := s.authRepo.ValidateRefreshToken(ctx, refreshToken)
	if err != nil || storedUserID != userIDStr {
		return nil, fmt.Errorf("token revoked or invalid")
	}

	s.authRepo.RevokeRefreshToken(ctx, refreshToken)

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user id format")
	}
	
	user, err := s.userRepo.GetByIDContext(ctx, uint(userID))
	if err != nil {
		return nil, err
	}

	return s.generateTokenPair(user.GetID(), user.Email, user.Role)
}

func (s *AuthService) ValidateAccessToken(tokenString string) (*entity.Claims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.accessSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid || claims["token_type"] != "access" {
		return nil, fmt.Errorf("invalid access token")
	}

	return &entity.Claims{
		UserID: claims["user_id"].(string),
		Email:  claims["email"].(string),
		Role:   claims["role"].(string),
	}, nil
}