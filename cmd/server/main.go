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
    "github.com/cureerel/gotemplate/internal/infrastructure"
    "github.com/cureerel/gotemplate/internal/infrastructure/dbtypes"
    "github.com/cureerel/gotemplate/internal/infrastructure/postgres/repositories"
    "github.com/cureerel/gotemplate/internal/interfaces/http/handler"
    "github.com/cureerel/gotemplate/internal/interfaces/http/router"
    "github.com/cureerel/gotemplate/pkg/logger"
    "gopkg.in/yaml.v3"
    "gorm.io/gorm"
)

type Config struct {
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
    JWT      JWTConfig      `yaml:"jwt"`
    Webhook  WebhookConfig  `yaml:"webhook"`
    Redis    RedisConfig    `yaml:"redis"`
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
    log := logger.New()

    // ── Config ────────────────────────────────────────────────
    _, b, _, _ := runtime.Caller(0)
    basepath   := filepath.Join(filepath.Dir(b), "../..")
    configPath := filepath.Join(basepath, "configs", "config.yaml")

    cfg, err := LoadConfig(configPath)
    if err != nil {
        log.Fatal("Failed to load config", logger.Field{Key: "error", Value: err})
    }

    if cfg.Server.Env == "production" {
        os.Setenv("GIN_MODE", "release")
    }

    // ── Database ──────────────────────────────────────────────
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

    // ── Repositories ──────────────────────────────────────────
    userRepo       := repositories.NewUserRepository(db)
    blogRepo       := repositories.NewBlogRepository(db)
    authRepo       := repositories.NewAuthRepository(db)
    sessionRepo    := repositories.NewSessionRepository(db) 
    webhookRepo    := repositories.NewWebhookRepository(db)
    productRepo    := repositories.NewProductRepository(db)
    orderRepo      := repositories.NewOrderRepository(db)
    membershipRepo := repositories.NewMembershipRepository(db)

    // ── Services ──────────────────────────────────────────────
    userService       := service.NewUserService(userRepo)
    blogService       := service.NewBlogService(blogRepo)
   authService := service.NewAuthService(userRepo, authRepo, sessionRepo, service.JWTConfig{
    AccessSecret:  cfg.JWT.AccessSecret,
    RefreshSecret: cfg.JWT.RefreshSecret,
})
    webhookService    := service.NewWebhookService(webhookRepo, service.WebhookConfig{
        StripeSecret:   cfg.Webhook.StripeSecret,
        RazorpaySecret: cfg.Webhook.RazorpaySecret,
    })
    productService    := service.NewProductService(productRepo)
    orderService      := service.NewOrderService(orderRepo, productRepo)
    paymentService    := service.NewPaymentService(webhookRepo, orderRepo)
    membershipService := service.NewMembershipService(membershipRepo)

    // ── Handlers ──────────────────────────────────────────────
    userHandler       := handler.NewUserHandler(userService)
    blogHandler       := handler.NewBlogHandler(blogService)
    authHandler       := handler.NewAuthHandler(authService)
    webhookHandler    := handler.NewWebhookHandler(webhookService, log)
    productHandler    := handler.NewProductHandler(productService)
    orderHandler      := handler.NewOrderHandler(orderService)
    paymentHandler    := handler.NewPaymentHandler(paymentService)
    membershipHandler := handler.NewMembershipHandler(membershipService)

    // ── Router ────────────────────────────────────────────────
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
    )

    // ── Server ────────────────────────────────────────────────
    srv := &http.Server{
        Addr:         ":" + cfg.Server.Port,
        Handler:      r,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    go func() {
        log.Info("Server starting", logger.Field{Key: "addr", Value: "0.0.0.0:" + cfg.Server.Port})
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal("Failed to run server", logger.Field{Key: "error", Value: err})
        }
    }()

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