package dto

import (
	"time"
)

// RegisterRequest defines user registration input
type RegisterRequest struct {
	Username string `json:"username" validate:"required,username"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	FullName string `json:"full_name" validate:"required,max=100"`
}

// LoginRequest defines user authentication input
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// TokenResponse defines authentication token output
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	// RefreshToken string    `json:"refresh_token"` unused
	ExpiresAt time.Time `json:"expires_at"`
	TokenType string    `json:"token_type"`
}

/**
/ RefreshRequest defines token refresh input
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
*/
