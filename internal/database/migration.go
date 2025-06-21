package database

import (
	"boreholedata-ms/internal/models"

	"gorm.io/gorm"
)

// RunMigrations executes the GORM auto-migration for all application models, centralize migration logic
func RunMigrations(db *gorm.DB) error {
	// GORM's AutoMigrate will create tables, add missing columns, and indices.
	return db.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.Station{},
		&models.LithologyLog{},
		&models.LabSample{},
		&models.UCSResult{},
		&models.UserActivity{},
	)
}