package database

import (
    "fmt"
    "log"
    "os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func NewPostgresConnection(dsn string) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.New(
            log.New(os.Stdout, "\n", log.LstdFlags),
            logger.Config{
                LogLevel: logger.Info,
                Colorful: true,
            },
        ),
    })

    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    return db, nil
}