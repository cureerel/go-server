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
	"github.com/cureerel/gotemplate/internal/infrastructure"
	"github.com/cureerel/gotemplate/internal/infrastructure/dbtypes"
	"github.com/cureerel/gotemplate/internal/infrastructure/postgres/repositories"
	"github.com/cureerel/gotemplate/internal/interfaces/http/handler"
	"github.com/cureerel/gotemplate/internal/interfaces/http/router"
	"github.com/cureerel/gotemplate/pkg/logger"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

// ---------------- Config Structs ----------------
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Webhook  WebhookConfig  `yaml:"webhook"`
	Redis    RedisConfig    `yaml:"redis"`
	CORS     CORSConfig     `yaml:"cors"`
}

type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Env  string `yaml:"env"`
}

type DatabaseConfig struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type WebhookConfig struct {
	StripeSecret   string `yaml:"stripe_secret"`
	RazorpaySecret string `yaml:"razorpay_secret"`
}

type JWTConfig struct {
	AccessSecret  string `yaml:"access_secret"`
	RefreshSecret string `yaml:"refresh_secret"`
}

// ---------------- Load Config ----------------
func LoadConfig() (*Config, error) {
	var cfg Config

	// Try loading config.yaml from a few conventional locations.
	candidates := []string{
		"configs/config.yaml",
		"/app/configs/config.yaml",
		"/etc/app/config.yaml",
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
			}
			if err := yaml.Unmarshal(data, &cfg); err != nil {
				return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
			}
			break
		}
	}

	// ENV vars always win for these fields
	if v := os.Getenv("PORT"); v != "" {
		cfg.Server.Port = v
	}
	if v := os.Getenv("ENV"); v != "" {
		cfg.Server.Env = v
	}
	if v := os.Getenv("DATABASE_URL"); v != "" {
		cfg.Database.DSN = v
	}
	if v := os.Getenv("JWT_ACCESS_SECRET"); v != "" {
		cfg.JWT.AccessSecret = v
	}
	if v := os.Getenv("JWT_REFRESH_SECRET"); v != "" {
		cfg.JWT.RefreshSecret = v
	}
	if v := os.Getenv("STRIPE_SECRET"); v != "" {
		cfg.Webhook.StripeSecret = v
	}
	if v := os.Getenv("RAZORPAY_SECRET"); v != "" {
		cfg.Webhook.RazorpaySecret = v
	}
	if v := os.Getenv("REDIS_ADDR"); v != "" {
		cfg.Redis.Addr = v
	}
	if v := os.Getenv("REDIS_USERNAME"); v != "" {
		cfg.Redis.Username = v
	}
	if v := os.Getenv("REDIS_PASSWORD"); v != "" {
		cfg.Redis.Password = v
	}

	// CORS: Merge config.yaml origins with env var origins (env appends to config)
	if v := os.Getenv("CORS_ALLOWED_ORIGINS"); v != "" {
		for _, origin := range splitAndTrim(v, ",") {
			if origin != "" && !contains(cfg.CORS.AllowedOrigins, origin) {
				cfg.CORS.AllowedOrigins = append(cfg.CORS.AllowedOrigins, origin)
			}
		}
	}

	return &cfg, nil
}

func splitAndTrim(s, sep string) []string {
	parts := []string{}
	for _, part := range strings.Split(s, sep) {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ---------------- Main ----------------
func main() {
	log := logger.New()

	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config", logger.Field{Key: "error", Value: err})
	}

	if cfg.Server.Env == "production" {
		os.Setenv("GIN_MODE", "release")
	}

	// ----- Database -----
	dbClient, err := infrastructure.NewDatabase(infrastructure.DBConfig{
		Driver: "postgres",
		DSN:    cfg.Database.DSN,
	})
	if err != nil {
		log.Fatal("Database connection failed", logger.Field{Key: "error", Value: err})
	}
	defer dbClient.Close()

	sqlDB, ok := dbClient.(dbtypes.SQLDB)
	if !ok {
		log.Fatal("Database client is not SQLDB")
	}
	db := sqlDB.GormDB().(*gorm.DB)

	// ----- Repositories -----
	userRepo := repositories.NewUserRepository(db)
	blogRepo := repositories.NewBlogRepository(db)
	authRepo := repositories.NewAuthRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	webhookRepo := repositories.NewWebhookRepository(db)
	productRepo := repositories.NewProductRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	membershipRepo := repositories.NewMembershipRepository(db)

	// ----- Services -----
	userService := service.NewUserService(userRepo)
	blogService := service.NewBlogService(blogRepo)
	authService := service.NewAuthService(userRepo, authRepo, sessionRepo, service.JWTConfig{
		AccessSecret:  cfg.JWT.AccessSecret,
		RefreshSecret: cfg.JWT.RefreshSecret,
	})
	webhookService := service.NewWebhookService(webhookRepo, service.WebhookConfig{
		StripeSecret:   cfg.Webhook.StripeSecret,
		RazorpaySecret: cfg.Webhook.RazorpaySecret,
	})
	productService := service.NewProductService(productRepo)
	orderService := service.NewOrderService(orderRepo, productRepo)
	paymentService := service.NewPaymentService(webhookRepo, orderRepo)
	membershipService := service.NewMembershipService(membershipRepo)

	// ----- Handlers -----
	userHandler := handler.NewUserHandler(userService)
	blogHandler := handler.NewBlogHandler(blogService)
	authHandler := handler.NewAuthHandler(authService)
	webhookHandler := handler.NewWebhookHandler(webhookService, log)
	productHandler := handler.NewProductHandler(productService)
	orderHandler := handler.NewOrderHandler(orderService)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	membershipHandler := handler.NewMembershipHandler(membershipService)

	// ----- Router -----
	r := router.SetupRouter(
		userHandler,
		blogHandler,
		authHandler,
		authService,
		webhookHandler,
		productHandler,
		orderHandler,
		paymentHandler,
		membershipHandler,
		log,
		cfg.CORS.AllowedOrigins,
	)

	// ----- Server -----
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("Server starting", logger.Field{Key: "addr", Value: "0.0.0.0:" + port})
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to run server", logger.Field{Key: "error", Value: err})
		}
	}()

	// ----- Graceful Shutdown -----
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", logger.Field{Key: "error", Value: err})
	}

	log.Info("Server exited")
}