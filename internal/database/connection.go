package database

import (
	"boreholedata-ms/config" // Changed import path to the new root-level config package
	"boreholedata-ms/internal/exception"
	"fmt"
	"log" // Import log for direct logging in this package

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger" // Import GORM logger
)

// InitDB initializes the GORM database connection.
// It takes a config.Config pointer and returns a *gorm.DB instance
// and a *exception.AppError if any error occurs during connection.
func InitDB(cfg *config.Config) (*gorm.DB, *exception.AppError) {
	// Construct the DSN (Data Source Name) for MySQL using values from the Config struct.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Open the GORM database connection.
	// IMPORTANT: Set logger.Info here to get detailed SQL logs and errors from GORM.
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Changed from logger.Silent to logger.Info
	})
	if err != nil {
		log.Printf("ERROR: GORM failed to open database connection: %v", err) // Log the underlying error
		return nil, exception.NewDatabaseError("DB connection failed", err)
	}
	log.Println("DEBUG: GORM database connection opened.")

	// Get the underlying sql.DB instance from GORM for pinging.
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("ERROR: Failed to get underlying SQL DB instance: %v", err)
		return nil, exception.NewDatabaseError("DB instance retrieval failed", err)
	}
	log.Println("DEBUG: Underlying SQL DB instance retrieved.")

	// Ping the database to verify the connection is alive.
	if err := sqlDB.Ping(); err != nil {
		log.Printf("ERROR: Database ping failed: %v", err) // Log the underlying error
		return nil, exception.NewDatabaseError("DB ping failed", err)
	}
	log.Println("DEBUG: Database ping successful.")

	// Return the GORM database instance and nil for the AppError on success.
	return db, nil
}

// MigrateSchema performs database auto-migration for all models.
// This function is provided for convenience during development.
// In a production environment, consider using a dedicated migration tool (e.g., golang-migrate).
// It's crucial to ensure all models are correctly defined before calling this.
// Note: This function is not called by default in your main.go,
// instead, db.AutoMigrate is called directly in main.go.
// If you want to use this, you'd replace the db.AutoMigrate block in main.go
// with a call to this function.
/*
func MigrateSchema(db *gorm.DB) error {
	log.Println("Attempting database auto-migration...")
	err := db.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.Station{},
		&models.LithologyLog{},
		&models.LabSample{},
		&models.UCSResult{},
	)
	if err != nil {
		log.Printf("ERROR: Database auto-migration failed: %v", err)
		return fmt.Errorf("database auto-migration failed: %w", err)
	}
	log.Println("Database auto-migration completed.")
	return nil
}
*/
