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

    "github.com/cureerel/gotemplate/internal/domain/entity"
    "github.com/cureerel/gotemplate/internal/domain/repository"
    "github.com/cureerel/gotemplate/pkg/apperror"
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    userRepo      repository.UserRepository
    authRepo      repository.AuthRepository
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

func NewAuthService(userRepo repository.UserRepository, authRepo repository.AuthRepository, cfg JWTConfig) *AuthService {
    return &AuthService{
        userRepo:      userRepo,
        authRepo:      authRepo,
        accessSecret:  []byte(cfg.AccessSecret),
        refreshSecret: []byte(cfg.RefreshSecret),
        tokenHashKey:  []byte(cfg.TokenHashKey),
        accessTTL:     15 * time.Minute,
        refreshTTL:    7 * 24 * time.Hour,
    }
}

// Signup creates a new user with hashed password
func (s *AuthService) Signup(ctx context.Context, name, email, password string) (*entity.User, error) {
    existing, err := s.userRepo.GetByEmail(ctx, email)
    if err != nil {
        return nil, apperror.NewInternal(err, "failed to check existing user")
    }
    if existing != nil {
        return nil, apperror.NewBadRequest("email already registered")
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, apperror.NewInternal(err, "failed to hash password")
    }

    user := &entity.User{
        Name:         name,
        Email:        email,
        PasswordHash: string(hash),
        Role:         "user",
        IsActive:     true,
    }

    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, apperror.NewInternal(err, "failed to create user")
    }

    return user, nil
}

// Login authenticates and returns access + refresh tokens
func (s *AuthService) Login(ctx context.Context, email, password string) (*entity.AuthToken, error) {
    user, err := s.userRepo.GetByEmail(ctx, email)
    if err != nil || user == nil {
        return nil, apperror.NewBadRequest("invalid credentials")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
        return nil, apperror.NewBadRequest("invalid credentials")
    }

    return s.generateTokenPair(ctx, user.ID, user.Email, user.Role)
}

// generateTokenPair creates access + refresh JWTs
func (s *AuthService) generateTokenPair(ctx context.Context, userID uint, email, role string) (*entity.AuthToken, error) {
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

    refreshClaims := jwt.MapClaims{
        "user_id":    strconv.Itoa(int(userID)),
        "token_type": "refresh",
        "exp":        now.Add(s.refreshTTL).Unix(),
        "iat":        now.Unix(),
        "jti":        uuid.New().String(),
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

    return &entity.AuthToken{
        AccessToken:  accessStr,
        RefreshToken: refreshStr,
        ExpiresAt:    now.Add(s.accessTTL),
    }, nil
}

// hashToken safely hashes token
func hashToken(token string, key []byte) string {
    h := hmac.New(sha256.New, key)
    h.Write([]byte(token))
    return hex.EncodeToString(h.Sum(nil))
}


// Refresh issues a new access + refresh token pair
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*entity.AuthToken, error) {
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

    if err := s.authRepo.RevokeRefreshToken(ctx, tokenHash); err != nil {
        return nil, apperror.NewInternal(err, "failed to revoke old refresh token")
    }

    user, err := s.userRepo.GetByID(ctx, stored.UserID)
    if err != nil || user == nil {
        return nil, apperror.NewBadRequest("user not found")
    }

    return s.generateTokenPair(ctx, user.ID, user.Email, user.Role)
}

// Logout revokes refresh token and blacklists access token
func (s *AuthService) Logout(ctx context.Context, accessToken, refreshToken string) error {
    if refreshToken != "" {
        tokenHash := hashToken(refreshToken, s.tokenHashKey)
        if err := s.authRepo.RevokeRefreshToken(ctx, tokenHash); err != nil {
            return apperror.NewInternal(err, "failed to revoke refresh token")
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