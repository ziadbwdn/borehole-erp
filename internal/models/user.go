package models

import (
	"time"

	"boreholedata-ms/internal/utils" // Import utils for BinaryUUID
)

// User represents a user in the system.
type User struct {
	ID           utils.BinaryUUID `gorm:"primaryKey;type:binary(16)" json:"id"` // Changed to utils.BinaryUUID
	Username     string           `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Email        string           `gorm:"size:100;not null;uniqueIndex" json:"email"`
	PasswordHash string           `gorm:"size:255;not null" json:"-"` // Stored hashed password, omitted from JSON
	FullName     string           `gorm:"size:100;not null" json:"fullName"`
	Role         UserRole         `gorm:"size:20;default:'geologist'" json:"role"`
	IsActive     bool             `gorm:"default:true" json:"isActive"`
	CreatedAt    time.Time        `gorm:"autoCreateTime" json:"createdAt"` // Corrected: Use autoCreateTime
	UpdatedAt    time.Time        `gorm:"autoUpdateTime" json:"updatedAt"` // Corrected: Use autoUpdateTime
}

// UserRole defines the type for user roles.
type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleGeologist UserRole = "geologist"
	RoleEngineer  UserRole = "engineer"
)
