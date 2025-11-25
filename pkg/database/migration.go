package database

import (
	"log"

	"github.com/brunobarlari/inventorypulse/internal/domain/models"
	"gorm.io/gorm"
)

// RunMigrations runs auto-migrations for all models
func RunMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
	)

	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

