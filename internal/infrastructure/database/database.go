package database

import (
	"fmt"

	"github.com/cureerel/gotemplate/internal/infrastructure/database/postgres"

	"gorm.io/gorm"
)

type Config struct {
	Driver string
	DSN    string
}

func New(cfg Config) (*gorm.DB, error) {
	switch cfg.Driver {

	case "postgres":
		return postgres.New(postgres.Config{
			DSN: cfg.DSN,
		})

	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}
}