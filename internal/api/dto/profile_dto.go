package dto

import (
	"boreholedata-ms/internal/utils"
	"time"
)

// ProfileResponse (assuming it exists and is used elsewhere for user profiles)
type ProfileResponse struct {
	ID        utils.BinaryUUID `json:"id"`
	Username  string           `json:"username"`
	Email     string           `json:"email"`
	FullName  string           `json:"fullName"`
	Role      string           `json:"role"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt time.Time        `json:"updatedAt"`
}

// UpdateProfileRequest (assuming it exists and is used elsewhere for user profile updates)
type UpdateProfileRequest struct {
	Email    *string `json:"email" validate:"omitempty,email"`
	FullName *string `json:"fullName" validate:"omitempty,max=100"`
	Password *string `json:"password" validate:"omitempty,password"`
}
