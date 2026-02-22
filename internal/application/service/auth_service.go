package service

import (
    "context"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
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

func (s *AuthService) Signup(ctx context.Context, name, email, password string) (*entity.User, error) {
    existing, err := s.userRepo.GetByEmail(ctx, email)
    if err != nil {
        return nil, err
    }
    if existing != nil {
        return nil, fmt.Errorf("email already registered")
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    user := &entity.User{
        Name:         name,
        Email:        email,
        PasswordHash: string(hashedPassword),
        Role:         "user",
        IsActive:     true,
    }

    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }
    return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*entity.AuthToken, error) {
    user, err := s.userRepo.GetByEmail(ctx, email)
    if err != nil {
        return nil, fmt.Errorf("invalid credentials")
    }
    if user == nil {
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

    return s.generateTokenPair(ctx, user.GetID(), user.Email, role)
}

func (s *AuthService) generateTokenPair(ctx context.Context, userID, email, role string) (*entity.AuthToken, error) {
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

    tokenHash := hashToken(refreshString)
    refreshEntity := &entity.RefreshToken{
        UserID:    parseUserID(userID),
        TokenHash: tokenHash,
        ExpiresAt: now.Add(s.refreshTTL),
    }
    if err := s.authRepo.SaveRefreshToken(ctx, refreshEntity); err != nil {
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
    tokenHash := hashToken(refreshToken)

    stored, err := s.authRepo.GetRefreshToken(ctx, tokenHash)
    if err != nil || stored == nil || stored.Revoked {
        return nil, fmt.Errorf("token revoked or invalid")
    }

    if fmt.Sprintf("%d", stored.UserID) != userIDStr {
        return nil, fmt.Errorf("token mismatch")
    }

    if err := s.authRepo.RevokeRefreshToken(ctx, tokenHash); err != nil {
        return nil, err
    }

    userID, err := strconv.ParseUint(userIDStr, 10, 64)
    if err != nil {
        return nil, fmt.Errorf("invalid user id format")
    }

    user, err := s.userRepo.GetByID(ctx, uint(userID))
    if err != nil || user == nil {
        return nil, fmt.Errorf("user not found")
    }

    return s.generateTokenPair(ctx, user.GetID(), user.Email, user.Role)
}

func (s *AuthService) Logout(ctx context.Context, accessToken, refreshToken string) error {
    tokenHash := hashToken(refreshToken)
    if err := s.authRepo.RevokeRefreshToken(ctx, tokenHash); err != nil {
        return err
    }
    return s.authRepo.BlacklistToken(ctx, hashToken(accessToken), time.Now().Add(s.accessTTL))
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

func hashToken(token string) string {
    h := hmac.New(sha256.New, []byte("token-hash-secret"))
    h.Write([]byte(token))
    return hex.EncodeToString(h.Sum(nil))
}

func parseUserID(id string) uint {
    v, _ := strconv.ParseUint(id, 10, 64)
    return uint(v)
}