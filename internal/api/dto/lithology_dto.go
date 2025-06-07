package dto

import (
	"time"

	"boreholedata-ms/internal/utils" // Import utils for BinaryUUID
)

// CreateLithologyLogRequest represents the request body for creating a new lithology log entry.
type CreateLithologyLogRequest struct {
	StationID          utils.BinaryUUID `json:"station_id" binding:"required"`
	DepthFrom          string           `json:"depth_from" binding:"required"`             // Changed to string
	DepthTo            string           `json:"depth_to" binding:"required"`               // Changed to string
	LithologyType      string           `json:"lithology_type" binding:"required,max=100"` // Changed field name
	RockColor          string           `json:"rock_color" binding:"max=50"`
	GrainSize          string           `json:"grain_size" binding:"max=50"`
	Texture            string           `json:"texture" binding:"max=100"`
	Structure          string           `json:"structure" binding:"max=100"`
	Hardness           string           `json:"hardness" binding:"max=50"`
	Weathering         string           `json:"weathering" binding:"max=50"`
	Fracturing         string           `json:"fracturing" binding:"max=100"`
	Description        string           `json:"description"`
	RQDPercentage      string           `json:"rqd_percentage"`      // Changed to string
	RecoveryPercentage string           `json:"recovery_percentage"` // Changed to string
	LoggedBy           string           `json:"logged_by" binding:"max=100"`
	LoggedDate         *time.Time       `json:"logged_date,omitempty"` // Changed to pointer
}

// UpdateLithologyLogRequest represents the request body for updating an existing lithology log entry.
// Fields are pointers to allow partial updates.
type UpdateLithologyLogRequest struct {
	StationID          *utils.BinaryUUID `json:"station_id,omitempty"`
	DepthFrom          *string           `json:"depth_from,omitempty"`                                 // Changed to string
	DepthTo            *string           `json:"depth_to,omitempty"`                                   // Changed to string
	LithologyType      *string           `json:"lithology_type,omitempty" binding:"omitempty,max=100"` // Changed field name
	RockColor          *string           `json:"rock_color,omitempty" binding:"omitempty,max=50"`
	GrainSize          *string           `json:"grain_size,omitempty" binding:"omitempty,max=50"`
	Texture            *string           `json:"texture,omitempty" binding:"omitempty,max=100"`
	Structure          *string           `json:"structure,omitempty" binding:"omitempty,max=100"`
	Hardness           *string           `json:"hardness,omitempty" binding:"omitempty,max=50"`
	Weathering         *string           `json:"weathering,omitempty" binding:"omitempty,max=50"`
	Fracturing         *string           `json:"fracturing,omitempty" binding:"omitempty,max=100"`
	Description        *string           `json:"description,omitempty"`
	RQDPercentage      *string           `json:"rqd_percentage,omitempty"`      // Changed to string
	RecoveryPercentage *string           `json:"recovery_percentage,omitempty"` // Changed to string
	LoggedBy           *string           `json:"logged_by,omitempty" binding:"omitempty,max=100"`
	LoggedDate         *time.Time        `json:"logged_date,omitempty"`
}

// LithologyLogResponse represents the response body for a lithology log entry.
type LithologyLogResponse struct {
	ID                 utils.BinaryUUID `json:"id"`
	StationID          utils.BinaryUUID `json:"station_id"`
	DepthFrom          string           `json:"depth_from"`     // Changed to string
	DepthTo            string           `json:"depth_to"`       // Changed to string
	LithologyType      string           `json:"lithology_type"` // Changed field name
	RockColor          string           `json:"rock_color,omitempty"`
	GrainSize          string           `json:"grain_size,omitempty"`
	Texture            string           `json:"texture,omitempty"`
	Structure          string           `json:"structure,omitempty"`
	Hardness           string           `json:"hardness,omitempty"`
	Weathering         string           `json:"weathering,omitempty"`
	Fracturing         string           `json:"fracturing,omitempty"`
	Description        string           `json:"description,omitempty"`
	RQDPercentage      string           `json:"rqd_percentage,omitempty"`
	RecoveryPercentage string           `json:"recovery_percentage,omitempty"`
	LoggedBy           string           `json:"logged_by,omitempty"`
	LoggedDate         *time.Time       `json:"logged_date,omitempty"` // Changed to pointer
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at,omitempty"`
}

// ListLithologyLogsResponse for listing multiple lithology logs.
type ListLithologyLogsResponse struct {
	Logs  []LithologyLogResponse `json:"logs"`
	Total int                    `json:"total"`
}
