package dto

import (
	"time"
)

// RegisterRequest defines user registration input
type RegisterRequest struct {
	Username string `json:"username" validate:"required,username"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	FullName string `json:"fullName" validate:"required,max=100"`
	Role     string `json:"role" validate:"required,max=100"` // User provides role during registration (will need validation)
}

// LoginRequest defines user authentication input
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// TokenResponse defines authentication token output
type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"` // Added refresh token
	ExpiresAt    time.Time `json:"expires_at"`    // Access token expiry
	TokenType    string    `json:"token_type"`
}

// ForgotPasswordRequest defines input for sending a password reset email.
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest defines input for resetting a password with a token.
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

// RefreshRequest defines token refresh input
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshResponse defines output for a token refresh request
type RefreshResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"` // New refresh token (for rotation)
	ExpiresAt    time.Time `json:"expires_at"`    // New access token expiry
	TokenType    string    `json:"token_type"`
}
