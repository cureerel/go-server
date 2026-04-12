// cmd/server/main.go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/cureerel/cserver/internal/application/service"
	"github.com/joho/godotenv"
	emailinfra "github.com/cureerel/cserver/internal/infrastructure/email"
	"github.com/cureerel/cserver/internal/infrastructure/email/resend"
	"github.com/cureerel/cserver/internal/infrastructure/postgres"
	"github.com/cureerel/cserver/internal/infrastructure/postgres/repositories"
	storageinfra "github.com/cureerel/cserver/internal/infrastructure/storage"
	"github.com/cureerel/cserver/internal/infrastructure/storage/cloudinary"
	"github.com/cureerel/cserver/internal/interfaces/http/handler"
	"github.com/cureerel/cserver/internal/interfaces/http/router"
	"github.com/cureerel/cserver/pkg/logger"
	"github.com/cureerel/cserver/internal/infrastructure/postgres/models"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	CORS     CORSConfig     `yaml:"cors"`
	Email    EmailConfig    `yaml:"email"`
	Storage  StorageConfig  `yaml:"storage"`
	Platform PlatformConfig `yaml:"platform"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Env  string `yaml:"env"`
}

type DatabaseConfig struct {
	DSN string `yaml:"dsn"`
}

type JWTConfig struct {
	AccessSecret  string `yaml:"access_secret"`
	RefreshSecret string `yaml:"refresh_secret"`
}

type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
}

type EmailConfig struct {
	ResendAPIKey string `yaml:"resend_api_key"`
	FromName     string `yaml:"from_name"`
	FromAddress  string `yaml:"from_address"`
}

type StorageConfig struct {
	CloudinaryCloudName string `yaml:"cloudinary_cloud_name"`
	CloudinaryAPIKey    string `yaml:"cloudinary_api_key"`
	CloudinaryAPISecret string `yaml:"cloudinary_api_secret"`
}

type PlatformConfig struct {
	OTPExpiryMinutes int `yaml:"otp_expiry_minutes"`
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load(".env")

	var cfg Config

	for _, path := range []string{"configs/config.yaml", "/app/configs/config.yaml"} {
		if _, err := os.Stat(path); err == nil {
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("failed to read config file: %w", err)
			}
			if err := yaml.Unmarshal(data, &cfg); err != nil {
				return nil, fmt.Errorf("failed to parse config file: %w", err)
			}
			break
		}
	}

	env := func(k, fb string) string {
		if v := os.Getenv(k); v != "" {
			fmt.Printf("[config] %s loaded from ENV\n", k)
			return v
		}
		if fb != "" {
			fmt.Printf("[config] %s loaded from YAML or default\n", k)
		}
		return fb
	}

	cfg.Server.Port = env("PORT", cfg.Server.Port)
	cfg.Server.Env = env("APP_ENV", cfg.Server.Env)
	cfg.Database.DSN = env("DATABASE_URL", env("EXTERNAL_DATABASE_URL", cfg.Database.DSN))
	cfg.JWT.AccessSecret = env("JWT_ACCESS_SECRET", cfg.JWT.AccessSecret)
	cfg.JWT.RefreshSecret = env("JWT_REFRESH_SECRET", cfg.JWT.RefreshSecret)
	cfg.Email.ResendAPIKey = env("RESEND_API_KEY", cfg.Email.ResendAPIKey)
	cfg.Email.FromName = env("EMAIL_FROM_NAME", cfg.Email.FromName)
	cfg.Email.FromAddress = env("EMAIL_FROM_ADDRESS", cfg.Email.FromAddress)
	cfg.Storage.CloudinaryCloudName = env("CLOUDINARY_CLOUD_NAME", cfg.Storage.CloudinaryCloudName)
	cfg.Storage.CloudinaryAPIKey = env("CLOUDINARY_API_KEY", cfg.Storage.CloudinaryAPIKey)
	cfg.Storage.CloudinaryAPISecret = env("CLOUDINARY_API_SECRET", cfg.Storage.CloudinaryAPISecret)

	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
		fmt.Printf("[config] PORT defaulting to 8080\n")
	}
	if cfg.Server.Env == "" {
		cfg.Server.Env = "development"
		fmt.Printf("[config] APP_ENV defaulting to development\n")
	}
	if cfg.Email.FromName == "" {
		cfg.Email.FromName = "Cureerel"
	}
	if cfg.Platform.OTPExpiryMinutes == 0 {
		cfg.Platform.OTPExpiryMinutes = 15
	}
	if v := os.Getenv("OTP_EXPIRY_MINUTES"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.Platform.OTPExpiryMinutes)
	}
	if v := os.Getenv("CORS_ALLOWED_ORIGINS"); v != "" {
		cfg.CORS.AllowedOrigins = nil
		for _, o := range strings.Split(v, ",") {
			if o = strings.TrimSpace(o); o != "" {
				cfg.CORS.AllowedOrigins = append(cfg.CORS.AllowedOrigins, o)
			}
		}
	}

	if cfg.Database.DSN == "" {
		return nil, fmt.Errorf("database DSN is required (set DATABASE_URL or EXTERNAL_DATABASE_URL env var, or database.dsn in config.yaml)")
	}

	return &cfg, nil
}

func main() {
	log := logger.New()

	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal("config load failed", logger.Field{Key: "error", Value: err})
	}

	if cfg.Server.Env == "production" {
		os.Setenv("GIN_MODE", "release")
	}

	dbClient, err := postgres.New(cfg.Database.DSN)
	if err != nil {
		log.Fatal("database failed", logger.Field{Key: "error", Value: err})
	}
	defer dbClient.Close()
	db := dbClient.GormDB().(*gorm.DB)

	// Inject auto-migration for newly added coin tables
	// to prevent 500 errors if Atlas is not installed.
	err = db.AutoMigrate(
		&models.UserWallet{}, 
		&models.CoinLedger{},
		&models.BlogUnlock{},
	)
	if err != nil {
		log.Error("AutoMigrate failed", logger.Field{Key: "error", Value: err})
	}

	var emailClient emailinfra.Provider
	if cfg.Email.ResendAPIKey != "" {
		emailClient = resend.New(cfg.Email.ResendAPIKey)
	} else {
		emailClient = &noopEmail{}
		log.Info("email client using noop (set RESEND_API_KEY to enable)")
	}

	var storageClient storageinfra.Provider
	if cfg.Storage.CloudinaryCloudName != "" && cfg.Storage.CloudinaryAPIKey != "" && cfg.Storage.CloudinaryAPISecret != "" {
		storageClient, err = cloudinary.New(
			cfg.Storage.CloudinaryCloudName,
			cfg.Storage.CloudinaryAPIKey,
			cfg.Storage.CloudinaryAPISecret,
		)
		if err != nil {
			log.Fatal("cloudinary init failed", logger.Field{Key: "error", Value: err})
		}
	} else {
		storageClient = &noopStorage{}
		log.Info("storage client using noop (set CLOUDINARY_* vars to enable)")
	}

	// ── Repositories ──────────────────────────────────────────
	userRepo         := repositories.NewUserRepository(db)
	blogRepo         := repositories.NewBlogRepository(db)
	authRepo         := repositories.NewAuthRepository(db)
	sessionRepo      := repositories.NewSessionRepository(db)
	otpRepo          := repositories.NewOTPRepository(db)
	serviceRepo      := repositories.NewServiceRepository(db)
	orderRepo        := repositories.NewOrderRepository(db)
	paymentRepo      := repositories.NewPaymentRepository(db)
	couponRepo       := repositories.NewCouponRepository(db)
	couponUsageRepo  := repositories.NewCouponUsageRepository(db)
	payoutRepo       := repositories.NewPayoutRepository(db)
	ticketRepo       := repositories.NewTicketRepository(db)
	upgradeRepo      := repositories.NewUpgradeRequestRepository(db)
	membershipRepo   := repositories.NewMembershipRepository(db)
	coinRepo         := repositories.NewCoinRepository(db)
	productRepo      := repositories.NewProductRepository(db)
	webhookRepo      := repositories.NewWebhookRepository(db)

	// ── Services ──────────────────────────────────────────────
	authService       := service.NewAuthService(userRepo, authRepo, sessionRepo, service.JWTConfig{
		AccessSecret: cfg.JWT.AccessSecret, RefreshSecret: cfg.JWT.RefreshSecret,
	})
	otpService        := service.NewOTPService(otpRepo, userRepo, emailClient, cfg.Email.FromName, cfg.Email.FromAddress, cfg.Platform.OTPExpiryMinutes)
	userService       := service.NewUserService(userRepo)
	blogService       := service.NewBlogService(blogRepo, coinRepo, membershipRepo)
	membershipService := service.NewMembershipService(membershipRepo)
	coinService       := service.NewCoinService(db, coinRepo, membershipService)
	serviceService    := service.NewServiceService(serviceRepo)
	orderService      := service.NewOrderService(db, orderRepo, serviceRepo, couponRepo, coinRepo)
	paymentService    := service.NewPaymentService(paymentRepo, orderRepo)
	couponService     := service.NewCouponService(couponRepo, couponUsageRepo, payoutRepo)
	payoutService     := service.NewPayoutService(payoutRepo)
	ticketService     := service.NewTicketService(ticketRepo)
	dashboardService  := service.NewDashboardService(db)
	superAdminService := service.NewSuperAdminService(userRepo, upgradeRepo, db)
	productService    := service.NewProductService(productRepo)

	webhookService    := service.NewWebhookService(webhookRepo, service.WebhookConfig{
		StripeSecret:   os.Getenv("STRIPE_WEBHOOK_SECRET"),
		RazorpaySecret: os.Getenv("RAZORPAY_WEBHOOK_SECRET"),
	})

	// ── Handlers ──────────────────────────────────────────────
	authHandler       := handler.NewAuthHandler(authService, userService, otpService, cfg.Platform.OTPExpiryMinutes)
	userHandler       := handler.NewUserHandler(userService)
	blogHandler       := handler.NewBlogHandler(blogService, coinService)
	serviceHandler    := handler.NewServiceHandler(serviceService)
	orderHandler      := handler.NewOrderHandler(orderService)
	paymentHandler    := handler.NewPaymentHandler(paymentService)
	couponHandler     := handler.NewCouponHandler(couponService)
	payoutHandler     := handler.NewPayoutHandler(payoutService)
	ticketHandler     := handler.NewTicketHandler(ticketService)
	dashboardHandler  := handler.NewDashboardHandler(dashboardService)
	superadminHandler := handler.NewSuperAdminHandler(superAdminService)
	uploadHandler     := handler.NewUploadHandler(storageClient)
	membershipHandler := handler.NewMembershipHandler(membershipService)
	pgHandler         := handler.NewPaymentGatewayHandler(membershipService, coinService)
	coinHandler       := handler.NewCoinHandler(coinService)
	productHandler    := handler.NewProductHandler(productService)
	webhookHandler    := handler.NewWebhookHandler(webhookService, log)

	r := router.SetupRouter(
		userHandler, blogHandler, authHandler, authService,
		serviceHandler, orderHandler, paymentHandler,
		couponHandler, payoutHandler, ticketHandler,
		dashboardHandler, superadminHandler, uploadHandler,
		membershipHandler, pgHandler, coinHandler,
		productHandler, webhookHandler,
		log, cfg.CORS.AllowedOrigins,
	)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("server starting", logger.Field{Key: "addr", Value: ":" + cfg.Server.Port})
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error", logger.Field{Key: "error", Value: err})
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("forced shutdown", logger.Field{Key: "error", Value: err})
	}
	log.Info("server exited")
}

type noopEmail struct{}

func (n *noopEmail) Send(_ context.Context, e emailinfra.Email) error {
	fmt.Printf("[email] To: %v Subject: %s\n", e.To, e.Subject)
	return nil
}

type noopStorage struct{}

func (n *noopStorage) Upload(_ context.Context, _ storageinfra.UploadInput) (storageinfra.UploadResult, error) {
	return storageinfra.UploadResult{URL: "https://placeholder.com/image.jpg", Key: "noop"}, nil
}
func (n *noopStorage) Delete(_ context.Context, _ string) error { return nil }
func (n *noopStorage) SignedURL(_ context.Context, key string, _ time.Duration) (string, error) {
	return "https://placeholder.com/" + key, nil
}