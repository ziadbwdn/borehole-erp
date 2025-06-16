package models

import (
	"time"

	"boreholedata-ms/internal/utils" // Import utils for BinaryUUID
)

// User represents a user in the system.
type User struct {
	ID            utils.BinaryUUID `gorm:"primaryKey;type:binary(16)" json:"id"`
	Username      string           `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Email         string           `gorm:"size:100;not null;uniqueIndex" json:"email"`
	PasswordHash  string           `gorm:"size:255;not null" json:"-"`
	FullName      string           `gorm:"size:100;not null" json:"fullName"`
	Role          UserRole         `gorm:"size:20;default:'geologist'" json:"role"`
	IsActive      bool             `gorm:"default:true" json:"isActive"`
	EmailVerified bool             `gorm:"default:false" json:"emailVerified"`
	LastLoginAt   *time.Time       `json:"lastLoginAt,omitempty"`

	// --- Refresh Token Fields ---
	// Store refresh token hash (SHA256)
	RefreshToken          *string    `gorm:"size:64;uniqueIndex" json:"-"` // SHA256 hash is 64 chars (hex encoded)
	RefreshTokenExpiresAt *time.Time `json:"-"`
	// --- End Refresh Token Fields ---

	PasswordResetToken  *string    `json:"-"` // This will store SHA256 hash of reset token
	PasswordResetSentAt *time.Time `json:"-"`
	CreatedAt           time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt           time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}

type UserRole string

const (
	RoleAdmin         UserRole = "admin"
	RoleGeologist     UserRole = "geologist"
	RoleEngineer      UserRole = "engineer"
	RoleLabTechnician UserRole = "lab_technician"
	RoleGuest         UserRole = "guest"
)
