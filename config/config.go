package config

import (
	"fmt"
	"log" // Added log for debug prints
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application-specific configuration.
type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	JWTSecret  string
	Port       string
}

// LoadConfig loads configuration from environment variables or a .env file.
func LoadConfig() (*Config, error) {
	log.Println("Attempting to load .env file...") // Debug print
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found or could not be loaded (this is okay if variables are set directly): %v", err) // More descriptive warning
	} else {
		log.Println(".env file loaded successfully.") // Debug print
	}

	// --- NEW DEBUG LINE HERE ---
	rawDBPassword := os.Getenv("DB_PASSWORD")
	log.Printf("DEBUG: Raw DB_PASSWORD from os.Getenv: '%s' (length: %d)", rawDBPassword, len(rawDBPassword))
	// --- END NEW DEBUG LINE ---

	cfg := &Config{
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBName:     getEnv("DB_NAME", "boreholedata_db"),
		JWTSecret:  getEnv("JWT_SECRET", "supersecretjwtkey"),
		Port:       getEnv("PORT", "8080"),
	}

	// Debug print for loaded values
	log.Printf("Loaded DB_USER: %s", cfg.DBUser)
	log.Printf("Loaded DB_HOST: %s", cfg.DBHost)
	log.Printf("Loaded DB_PORT: %s", cfg.DBPort)
	log.Printf("Loaded DB_NAME: %s", cfg.DBName)
	// IMPORTANT: Do NOT log DB_PASSWORD directly in production. This is for debug only.
	log.Printf("Loaded DB_PASSWORD (length): %d", len(cfg.DBPassword))
	if len(cfg.DBPassword) > 0 {
		log.Println("DB_PASSWORD is NOT empty.")
	} else {
		log.Println("DB_PASSWORD IS EMPTY.")
	}
	log.Printf("Loaded JWT_SECRET (length): %d", len(cfg.JWTSecret))
	log.Printf("Loaded PORT: %s", cfg.Port)

	// Basic validation for critical configurations
	if cfg.DBUser == "" || cfg.DBHost == "" || cfg.DBName == "" || cfg.JWTSecret == "" {
		return nil, fmt.Errorf("critical database or JWT configuration missing. Please check DB_USER, DB_HOST, DB_NAME, JWT_SECRET")
	}

	// Optional: Validate port is a number
	if _, err := strconv.Atoi(cfg.Port); err != nil {
		return nil, fmt.Errorf("invalid PORT value: %v", err)
	}

	return cfg, nil
}

// getEnv retrieves an environment variable or returns a fallback value.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
