package dto

import (
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"time"
)

// CreateStationRequest represents the request body for creating a new station.
type CreateStationRequest struct {
	ProjectID      utils.BinaryUUID `json:"project_id" binding:"required"`
	StationCode    string           `json:"station_code" binding:"required,min=3,max=50"`
	StationName    string           `json:"station_name" binding:"max=100"`
	StationType    string           `json:"station_type" binding:"required,max=50"`
	Latitude       string           `json:"latitude" binding:"required"`
	Longitude      string           `json:"longitude" binding:"required"`
	Elevation      string           `json:"elevation" binding:"required"`
	GWL            string           `json:"groundwater_level" binding:"required"`
	TotalDepth     string           `json:"total_depth" binding:"required"`
	DrillingDate   *time.Time       `json:"drilling_date,omitempty"`
	DrillingStatus string           `json:"drilling_status,omitempty" binding:"omitempty,max=50"`
	GeologistName  string           `json:"geologist_name" binding:"max=100"`
	Notes          string           `json:"notes"`
}

// UpdateStationRequest represents the request body for updating an existing station.
// Fields are pointers to allow partial updates (nil means no change).
type UpdateStationRequest struct {
	StationCode    *string    `json:"station_code,omitempty" binding:"omitempty,min=3,max=50"`
	StationName    *string    `json:"station_name,omitempty" binding:"omitempty,max=100"`
	StationType    *string    `json:"station_type,omitempty" binding:"omitempty,max=50"`
	Latitude       *string    `json:"latitude,omitempty"`
	Longitude      *string    `json:"longitude,omitempty"`
	Elevation      *string    `json:"elevation,omitempty"`
	GWL            *string    `json:"groundwater_level,omitempty"`
	TotalDepth     *string    `json:"total_depth,omitempty"`
	DrillingDate   *time.Time `json:"drilling_date,omitempty"`
	DrillingStatus *string    `json:"drilling_status,omitempty" binding:"omitempty,max=50"` 
	GeologistName  *string    `json:"geologist_name,omitempty" binding:"omitempty,max=100"`
	Notes          *string    `json:"notes,omitempty"`
}

// StationResponse represents the response body for station details.
type StationResponse struct {
	ID             utils.BinaryUUID `json:"id"`
	ProjectID      utils.BinaryUUID `json:"project_id"`
	StationCode    string           `json:"station_code"`
	StationName    string           `json:"station_name"`
	StationType    string           `json:"station_type"`
	Latitude       string           `json:"latitude"`
	Longitude      string           `json:"longitude"`
	Elevation      string           `json:"elevation"`
	GWL            string           `json:"groundwater_level,omitempty"`
	TotalDepth     string           `json:"total_depth"`
	DrillingDate   *time.Time       `json:"drilling_date,omitempty"`
	DrillingStatus string           `json:"drilling_status"`
	GeologistName  string           `json:"geologist_name"`
	Notes          string           `json:"notes"`
	ProjectName    string           `json:"project_name,omitempty"` 
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at,omitempty"`
}

// StationPointGeoResponse: Geographic only
type StationPointGeoResponse struct {
	StationCode string `json:"station_code"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	Elevation   string `json:"elevation"`
}

// StationPointUTMResponse: Geographic and UTM  
type StationPointUTMResponse struct {
	StationCode string `json:"station_code"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	Elevation   string `json:"elevation"`
	UTMZone     string `json:"utm_zone"`
	Easting  string `json:"easting"`
	Northing string `json:"northing"`
}

// map station to response
func MapStationToResponse(station *models.Station) *StationResponse {
	if station == nil {
		return nil
	}

	return &StationResponse{
		ID:             station.ID,
		ProjectID:      station.ProjectID,
		StationCode:    station.StationCode,
		StationName:    station.StationName,
		StationType:    station.StationType,
		Latitude:       utils.GormDecimalToString(&station.Latitude),
		Longitude:      utils.GormDecimalToString(&station.Longitude),
		Elevation:      utils.GormDecimalToString(&station.Elevation),
		GWL:            utils.GormDecimalToString(&station.GWL),
		TotalDepth:     utils.GormDecimalToString(&station.TotalDepth),
		DrillingDate:   station.DrillingDate,
		DrillingStatus: string(station.DrillingStatus),
		GeologistName:  station.GeologistName,
		Notes:          station.Notes,
		CreatedAt:      station.CreatedAt,
		UpdatedAt:      station.UpdatedAt,
	}
}

// ListStationsResponse (optional, for listing multiple stations)
type ListStationsResponse struct {
	Stations []StationResponse `json:"stations"`
	Total    int               `json:"total"`
}
