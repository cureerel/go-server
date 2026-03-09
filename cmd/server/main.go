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

	"github.com/cureerel/gotemplate/internal/application/service"
	emailinfra "github.com/cureerel/gotemplate/internal/infrastructure/email"
	"github.com/cureerel/gotemplate/internal/infrastructure/email/resend"
	"github.com/cureerel/gotemplate/internal/infrastructure/postgres"
	"github.com/cureerel/gotemplate/internal/infrastructure/postgres/repositories"
	storageinfra "github.com/cureerel/gotemplate/internal/infrastructure/storage"
	"github.com/cureerel/gotemplate/internal/infrastructure/storage/cloudinary"
	"github.com/cureerel/gotemplate/internal/interfaces/http/handler"
	"github.com/cureerel/gotemplate/internal/interfaces/http/router"
	"github.com/cureerel/gotemplate/pkg/logger"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

// ── Config ────────────────────────────────────────────────────

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
	var cfg Config
	for _, path := range []string{"configs/config.yaml", "/app/configs/config.yaml"} {
		if _, err := os.Stat(path); err == nil {
			data, _ := os.ReadFile(path)
			_ = yaml.Unmarshal(data, &cfg)
			break
		}
	}

	env := func(key, fallback string) string {
		if v := os.Getenv(key); v != "" {
			return v
		}
		return fallback
	}

	cfg.Server.Port        = env("PORT",                 cfg.Server.Port)
	cfg.Server.Env         = env("APP_ENV",              cfg.Server.Env)
	cfg.Database.DSN       = env("DATABASE_URL",         cfg.Database.DSN)
	cfg.JWT.AccessSecret   = env("JWT_ACCESS_SECRET",    cfg.JWT.AccessSecret)
	cfg.JWT.RefreshSecret  = env("JWT_REFRESH_SECRET",   cfg.JWT.RefreshSecret)
	cfg.Email.ResendAPIKey = env("RESEND_API_KEY",       cfg.Email.ResendAPIKey)
	cfg.Email.FromName     = env("EMAIL_FROM_NAME",      cfg.Email.FromName)
	cfg.Email.FromAddress  = env("EMAIL_FROM_ADDRESS",   cfg.Email.FromAddress)
	cfg.Storage.CloudinaryCloudName = env("CLOUDINARY_CLOUD_NAME", cfg.Storage.CloudinaryCloudName)
	cfg.Storage.CloudinaryAPIKey    = env("CLOUDINARY_API_KEY",    cfg.Storage.CloudinaryAPIKey)
	cfg.Storage.CloudinaryAPISecret = env("CLOUDINARY_API_SECRET", cfg.Storage.CloudinaryAPISecret)

	if cfg.Server.Port == ""      { cfg.Server.Port = "8080" }
	if cfg.Email.FromName == ""   { cfg.Email.FromName = "Cureerel" }
	if cfg.Platform.OTPExpiryMinutes == 0 { cfg.Platform.OTPExpiryMinutes = 15 }
	if v := os.Getenv("OTP_EXPIRY_MINUTES"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.Platform.OTPExpiryMinutes)
	}

	if v := os.Getenv("CORS_ALLOWED_ORIGINS"); v != "" {
		for _, o := range strings.Split(v, ",") {
			if o = strings.TrimSpace(o); o != "" {
				cfg.CORS.AllowedOrigins = append(cfg.CORS.AllowedOrigins, o)
			}
		}
	}

	return &cfg, nil
}

// ── Main ──────────────────────────────────────────────────────

func main() {
	log := logger.New()

	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal("config load failed", logger.Field{Key: "error", Value: err})
	}
	if cfg.Server.Env == "production" {
		os.Setenv("GIN_MODE", "release")
	}

	// ── Database ──────────────────────────────────────────────
	// postgres.New returns (dbtypes.SQLDB, error).
	// GormDB() returns interface{} — assert to *gorm.DB for repositories.
	dbClient, err := postgres.New(cfg.Database.DSN)
	if err != nil {
		log.Fatal("database connection failed", logger.Field{Key: "error", Value: err})
	}
	defer dbClient.Close()

	db := dbClient.GormDB().(*gorm.DB)

	// ── Email ─────────────────────────────────────────────────
	var emailClient emailinfra.Provider
	if cfg.Email.ResendAPIKey != "" {
		emailClient = resend.New(cfg.Email.ResendAPIKey)
	} else {
		log.Info("RESEND_API_KEY not set — using noop email")
		emailClient = &noopEmail{}
	}

	// ── Storage ───────────────────────────────────────────────
	var storageClient storageinfra.Provider
	if cfg.Storage.CloudinaryCloudName != "" {
		storageClient = cloudinary.New(
			cfg.Storage.CloudinaryCloudName,
			cfg.Storage.CloudinaryAPIKey,
			cfg.Storage.CloudinaryAPISecret,
		)
	} else {
		log.Info("Cloudinary not configured — using noop storage")
		storageClient = &noopStorage{}
	}

	// ── Repositories ──────────────────────────────────────────
	userRepo    := repositories.NewUserRepository(db)
	blogRepo    := repositories.NewBlogRepository(db)
	authRepo    := repositories.NewAuthRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	otpRepo     := repositories.NewOTPRepository(db)
	serviceRepo := repositories.NewServiceRepository(db)

	// ── Services ──────────────────────────────────────────────
	authService := service.NewAuthService(userRepo, authRepo, sessionRepo, service.JWTConfig{
		AccessSecret:  cfg.JWT.AccessSecret,
		RefreshSecret: cfg.JWT.RefreshSecret,
	})
	otpService := service.NewOTPService(
		otpRepo, userRepo, emailClient,
		cfg.Email.FromName, cfg.Email.FromAddress,
		cfg.Platform.OTPExpiryMinutes,
	)
	userService    := service.NewUserService(userRepo)
	blogService    := service.NewBlogService(blogRepo)
	serviceService := service.NewServiceService(serviceRepo)

	// ── Handlers ──────────────────────────────────────────────
	authHandler    := handler.NewAuthHandler(authService, otpService, cfg.Platform.OTPExpiryMinutes)
	userHandler    := handler.NewUserHandler(userService)
	blogHandler    := handler.NewBlogHandler(blogService)
	serviceHandler := handler.NewServiceHandler(serviceService)
	uploadHandler  := handler.NewUploadHandler(storageClient)

	// ── Router ────────────────────────────────────────────────
	r := router.SetupRouter(
		userHandler,
		blogHandler,
		authHandler,
		authService,
		serviceHandler,
		uploadHandler,
		log,
		cfg.CORS.AllowedOrigins,
	)

	// ── HTTP server ───────────────────────────────────────────
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

// ── Noop providers ────────────────────────────────────────────

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