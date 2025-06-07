package dto

import (
	"boreholedata-ms/internal/utils"
	"time"
)

// CreateStationRequest represents the request body for creating a new station.
type CreateStationRequest struct {
	ProjectID     utils.BinaryUUID `json:"project_id" binding:"required"`
	StationCode   string           `json:"station_code" binding:"required,min=3,max=50"`
	StationName   string           `json:"station_name" binding:"max=100"`
	StationType   string           `json:"station_type" binding:"required,max=50"`
	Latitude      string           `json:"latitude" binding:"required"`    // Changed to string
	Longitude     string           `json:"longitude" binding:"required"`   // Changed to string
	Elevation     string           `json:"elevation" binding:"required"`   // Changed to string
	TotalDepth    string           `json:"total_depth" binding:"required"` // Changed to string
	DrillingDate    *time.Time       `json:"drilling_date,omitempty"` // Changed to pointer
	GeologistName string           `json:"geologist_name" binding:"max=100"`
	Notes         string           `json:"notes"`
}

// UpdateStationRequest represents the request body for updating an existing station.
// Fields are pointers to allow partial updates (nil means no change).
type UpdateStationRequest struct {
	StationCode   *string    `json:"station_code,omitempty" binding:"omitempty,min=3,max=50"`
	StationName   *string    `json:"station_name,omitempty" binding:"omitempty,max=100"`
	StationType   *string    `json:"station_type,omitempty" binding:"omitempty,max=50"`
	Latitude      *string    `json:"latitude,omitempty"`    // Changed to string
	Longitude     *string    `json:"longitude,omitempty"`   // Changed to string
	Elevation     *string    `json:"elevation,omitempty"`   // Changed to string
	TotalDepth    *string    `json:"total_depth,omitempty"` // Changed to string
	DrillingDate  *time.Time `json:"drilling_date,omitempty"`
	GeologistName *string    `json:"geologist_name,omitempty" binding:"omitempty,max=100"`
	Notes         *string    `json:"notes,omitempty"`
}

// StationResponse represents the response body for station details.
type StationResponse struct {
	ID            utils.BinaryUUID `json:"id"`
	ProjectID     utils.BinaryUUID `json:"project_id"`
	StationCode   string           `json:"station_code"`
	StationName   string           `json:"station_name"`
	StationType   string           `json:"station_type"`
	Latitude      string           `json:"latitude"`    // Changed to string
	Longitude     string           `json:"longitude"`   // Changed to string
	Elevation     string           `json:"elevation"`   // Changed to string
	TotalDepth    string           `json:"total_depth"` // Changed to string
	DrillingDate  *time.Time       `json:"drilling_date,omitempty"` // Changed to pointer
	GeologistName string           `json:"geologist_name"`
	Notes         string           `json:"notes"`
	ProjectName   string           `json:"project_name,omitempty"` // Added for response enrichment
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at,omitempty"`
}

// ListStationsResponse (optional, for listing multiple stations)
type ListStationsResponse struct {
	Stations []StationResponse `json:"stations"`
	Total    int               `json:"total"`
}
