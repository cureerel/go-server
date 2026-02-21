package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/cureerel/gotemplate/internal/application/service"
	"github.com/cureerel/gotemplate/internal/infrastructure/database"
	"github.com/cureerel/gotemplate/internal/infrastructure/persistence"
	"github.com/cureerel/gotemplate/internal/interfaces/http/handler"
	"github.com/cureerel/gotemplate/internal/interfaces/http/router"
	"github.com/cureerel/gotemplate/pkg/logger"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
  JWT      JWTConfig      `yaml:"jwt"`      // ADD
	Webhook  WebhookConfig  `yaml:"webhook"`  // ADD
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	DSN string `yaml:"dsn"`
}

type WebhookConfig struct {
	StripeSecret   string `yaml:"stripe_secret"`
	RazorpaySecret string `yaml:"razorpay_secret"`
}

type JWTConfig struct {
	AccessSecret  string `yaml:"access_secret"`
	RefreshSecret string `yaml:"refresh_secret"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	return &cfg, nil
}

func main() {
	// Init logger
	log := logger.New()

	// Load config
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Join(filepath.Dir(b), "../..")
	configPath := filepath.Join(basepath, "configs", "config.yaml")

	cfg, err := LoadConfig(configPath)
	if err != nil {
		log.Fatal("Failed to load config", logger.Field{Key: "error", Value: err})
	}

	// Connect DB
	db, err := database.New(database.Config{
		Driver: "postgres",
		DSN:    cfg.Database.DSN,
	})
	if err != nil {
		log.Fatal("Database connection failed", logger.Field{Key: "error", Value: err})
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get sql.DB", logger.Field{Key: "error", Value: err})
	}
	defer sqlDB.Close()

	// Init repositories
	userRepo := persistence.NewUserRepository(db)
	blogRepo := persistence.NewBlogRepository(db)
	authRepo := persistence.NewAuthRepository(db) // TODO: Add Redis for token blacklist
	webhookRepo := persistence.NewWebhookRepository(db) // ADD THIS

	// Init services
	userService := service.NewUserService(userRepo)
	blogService := service.NewBlogService(blogRepo)
	authService := service.NewAuthService(userRepo, authRepo, service.JWTConfig{
	AccessSecret:  cfg.JWT.AccessSecret,
	RefreshSecret: cfg.JWT.RefreshSecret,
})
	webhookService := service.NewWebhookService(webhookRepo, service.WebhookConfig{
	StripeSecret:   cfg.Webhook.StripeSecret,
	RazorpaySecret: cfg.Webhook.RazorpaySecret,
})

	// Init handlers
	userHandler := handler.NewUserHandler(userService)
	blogHandler := handler.NewBlogHandler(blogService)
	authHandler := handler.NewAuthHandler(authService)
	webhookHandler := handler.NewWebhookHandler(webhookService, log)

	// Setup router - ADD log parameter
	r := router.SetupRouter(userHandler, blogHandler, authHandler, authService, webhookHandler, log)

	// Graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	go func() {
		addr := "0.0.0.0:" + cfg.Server.Port
		log.Info("Server starting", logger.Field{Key: "addr", Value: addr})
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to run server", logger.Field{Key: "error", Value: err})
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", logger.Field{Key: "error", Value: err})
	}
	log.Info("Server exited")
}