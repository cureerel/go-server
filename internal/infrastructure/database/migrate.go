package database

import (
	"github.com/cureerel/gotemplate/internal/domain/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
		&entity.Blog{},
	)
}