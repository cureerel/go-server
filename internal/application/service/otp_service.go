// internal/application/service/otp_service.go
package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/cureerel/gotemplate/internal/domain/entity"
	"github.com/cureerel/gotemplate/internal/domain/repository"
	emailinfra "github.com/cureerel/gotemplate/internal/infrastructure/email"
	"github.com/cureerel/gotemplate/pkg/apperror"
	"golang.org/x/crypto/bcrypt"
)

type OTPService struct {
	otpRepo       repository.OTPRepository
	userRepo      repository.UserRepository
	emailProvider emailinfra.Provider
	fromName      string
	fromAddr      string
	expiryMins    int
}

func NewOTPService(
	otpRepo repository.OTPRepository,
	userRepo repository.UserRepository,
	emailProvider emailinfra.Provider,
	fromName, fromAddr string,
	expiryMins int,
) *OTPService {
	return &OTPService{
		otpRepo: otpRepo, userRepo: userRepo, emailProvider: emailProvider,
		fromName: fromName, fromAddr: fromAddr, expiryMins: expiryMins,
	}
}

func (s *OTPService) ExpiryMinutes() int { return s.expiryMins }

// SendRegisterOTP generates and emails an OTP for new registration.
func (s *OTPService) SendRegisterOTP(ctx context.Context, email string) error {
	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return apperror.NewInternal(err, "failed to check email")
	}
	if existing != nil {
		return apperror.NewBadRequest("email already registered")
	}
	return s.sendOTP(ctx, nil, email, entity.OTPTypeRegister)
}

// SendResetOTP generates and emails an OTP for password reset.
// Always returns nil (silent for enumeration protection).
func (s *OTPService) SendResetOTP(ctx context.Context, email string) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return nil // silent
	}
	return s.sendOTP(ctx, &user.ID, email, entity.OTPTypeReset)
}

// VerifyAndCreateUser verifies OTP then creates a verified user. Used by RegisterVerify handler.
func (s *OTPService) VerifyAndCreateUser(ctx context.Context, name, email, password, code string) (*entity.User, error) {
	if err := s.verifyOTP(ctx, email, code, entity.OTPTypeRegister); err != nil {
		return nil, err
	}
	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to check email")
	}
	if existing != nil {
		return nil, apperror.NewBadRequest("email already registered")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to hash password")
	}
	user := &entity.User{
		Name: name, Email: email, PasswordHash: string(hash),
		Role: entity.RoleUser, IsActive: true, IsVerified: true,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, apperror.NewInternal(err, "failed to create user")
	}
	return user, nil
}

// VerifyResetAndChangePassword verifies OTP then updates password. Used by PasswordResetVerify handler.
func (s *OTPService) VerifyResetAndChangePassword(ctx context.Context, email, code, newPassword string) error {
	if err := s.verifyOTP(ctx, email, code, entity.OTPTypeReset); err != nil {
		return err
	}
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return apperror.NewNotFound("user not found")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperror.NewInternal(err, "failed to hash password")
	}
	return s.userRepo.UpdatePassword(ctx, user.ID, string(hash))
}

// ── Private helpers ───────────────────────────────────────────

func (s *OTPService) sendOTP(ctx context.Context, userID *uint, email, otpType string) error {
	_ = s.otpRepo.DeleteExpiredByEmail(ctx, email, otpType)
	code, err := generate6Digit()
	if err != nil {
		return apperror.NewInternal(err, "failed to generate OTP")
	}
	otp := &entity.OTP{
		UserID: userID, Email: email, Code: code, Type: otpType,
		ExpiresAt: time.Now().Add(time.Duration(s.expiryMins) * time.Minute),
	}
	if err := s.otpRepo.Create(ctx, otp); err != nil {
		return apperror.NewInternal(err, "failed to save OTP")
	}
	subject, html := s.emailContent(otpType, code)
	return s.emailProvider.Send(ctx, emailinfra.Email{
		From:    fmt.Sprintf("%s <%s>", s.fromName, s.fromAddr),
		To:      []string{email},
		Subject: subject,
		HTML:    html,
	})
}

func (s *OTPService) verifyOTP(ctx context.Context, email, code, otpType string) error {
	otp, err := s.otpRepo.GetLatestByEmail(ctx, email, otpType)
	if err != nil {
		return apperror.NewInternal(err, "failed to fetch OTP")
	}
	if otp == nil || !otp.IsValid() {
		return apperror.NewBadRequest("OTP not found or expired")
	}
	if otp.Code != code {
		return apperror.NewBadRequest("invalid OTP")
	}
	return s.otpRepo.MarkUsed(ctx, otp.ID)
}

func (s *OTPService) emailContent(otpType, code string) (subject, html string) {
	expiry := s.expiryMins
	switch otpType {
	case entity.OTPTypeRegister:
		subject = "Your Cureerel verification code"
		html = fmt.Sprintf(`<div style="font-family:sans-serif;max-width:480px;margin:0 auto;padding:32px"><h2 style="color:#111">Verify your email</h2><p style="color:#555">Your one-time code expires in %d minutes.</p><div style="background:#f5f5f5;border-radius:12px;padding:24px 32px;text-align:center;margin:24px 0"><span style="font-size:40px;font-weight:800;letter-spacing:10px;color:#111">%s</span></div><p style="color:#999;font-size:12px">If you didn't request this, ignore this email.</p></div>`, expiry, code)
	case entity.OTPTypeReset:
		subject = "Your Cureerel password reset code"
		html = fmt.Sprintf(`<div style="font-family:sans-serif;max-width:480px;margin:0 auto;padding:32px"><h2 style="color:#111">Reset your password</h2><p style="color:#555">Use this code within %d minutes.</p><div style="background:#f5f5f5;border-radius:12px;padding:24px 32px;text-align:center;margin:24px 0"><span style="font-size:40px;font-weight:800;letter-spacing:10px;color:#111">%s</span></div><p style="color:#999;font-size:12px">If you didn't request a password reset, you can safely ignore this email.</p></div>`, expiry, code)
	}
	return
}

func generate6Digit() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

// VerifyRegisterOTP is an alias used by the auth handler.
func (s *OTPService) VerifyRegisterOTP(ctx context.Context, email, code string) error {
	return s.verifyOTP(ctx, email, code, entity.OTPTypeRegister)
}

// VerifyResetOTP verifies the reset OTP and changes the password.
func (s *OTPService) VerifyResetOTP(ctx context.Context, email, code, newPassword string) error {
	return s.VerifyResetAndChangePassword(ctx, email, code, newPassword)
}