// internal/infrastructure/postgres/postgres.go
package postgres

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cureerel/gotemplate/internal/infrastructure/dbtypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Client struct {
	db *gorm.DB
}

var _ dbtypes.SQLDB = (*Client)(nil)

func New(dsn string) (dbtypes.SQLDB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("postgres connection failed: %w", err)
	}

	return &Client{db: db}, nil
}

func (c *Client) GormDB() interface{} { return c.db }

func (c *Client) Ping(ctx context.Context) error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

func (c *Client) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}