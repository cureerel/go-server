package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/cureerel/gotemplate/internal/application/service"
	"github.com/cureerel/gotemplate/internal/domain/entity"
	"github.com/cureerel/gotemplate/internal/infrastructure/database"
	"github.com/cureerel/gotemplate/internal/infrastructure/persistence"
	"github.com/cureerel/gotemplate/internal/interfaces/http/handler"
	"github.com/cureerel/gotemplate/internal/interfaces/http/router"
	"gopkg.in/yaml.v3"
)

// Config structs
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	DSN string `yaml:"dsn"`
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
	// Find project root
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Join(filepath.Dir(b), "../..")
	configPath := filepath.Join(basepath, "configs", "config.yaml")

	// Load config
	cfg, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect DB
	db, err := database.NewPostgresConnection(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate all entities
	if err := db.AutoMigrate(&entity.User{}, &entity.Blog{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// ==================== USER LAYERS ====================
	userRepo := persistence.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// ==================== BLOG LAYERS (NEW) ====================
	blogRepo := persistence.NewBlogRepository(db)
	blogService := service.NewBlogService(blogRepo)
	blogHandler := handler.NewBlogHandler(blogService)

	// Setup Router with both handlers
	r := router.SetupRouter(userHandler, blogHandler)

	// Start server
	addr := "0.0.0.0:" + cfg.Server.Port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}