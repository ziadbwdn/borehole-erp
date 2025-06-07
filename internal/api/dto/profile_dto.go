package dto

import (
	"boreholedata-ms/internal/utils"
	"time"
)

// ProfileResponse defines user profile output
type ProfileResponse struct {
	ID        utils.BinaryUUID `json:"id"`
	Username  string           `json:"username"`
	Email     string           `json:"email"`
	FullName  string           `json:"full_name"`
	Role      string           `json:"role"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt time.Time        `json:"updatedAt,omitempty"`
}

// UpdateProfileRequest represents the request body for updating a user's profile.
// Fields are pointers to allow partial updates (nil means no change).
type UpdateProfileRequest struct {
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
	FullName *string `json:"fullName,omitempty" binding:"omitempty,min=3,max=100"`
	Password *string `json:"password,omitempty" binding:"omitempty,min=12"` // Password validation handled in service
}
