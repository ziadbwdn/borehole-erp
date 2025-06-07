package dto

import (
	"boreholedata-ms/internal/utils"
	"time"
)

// CreateUCSResultRequest represents the request body for creating a new UCS test result.
type CreateUCSResultRequest struct {
	SampleID         utils.BinaryUUID `json:"sample_id" binding:"required"`
	UCSValue         string           `json:"ucs_value" binding:"required"` // Changed to string
	Unit             string           `json:"unit" binding:"max=10"`
	TestMethod       string           `json:"test_method" binding:"max=50"`
	SpecimenDiameter string           `json:"specimen_diameter"` // Changed to string
	SpecimenHeight   string           `json:"specimen_height"`   // Changed to string
	FailureMode      string           `json:"failure_mode" binding:"max=100"`
	Notes            string           `json:"notes"`
}

// UpdateUCSResultRequest represents the request body for updating an existing UCS test result.
type UpdateUCSResultRequest struct {
	SampleID         *utils.BinaryUUID `json:"sample_id,omitempty"`
	UCSValue         *string           `json:"ucs_value,omitempty"` // Changed to string
	Unit             *string           `json:"unit,omitempty" binding:"omitempty,max=10"`
	TestMethod       *string           `json:"test_method,omitempty" binding:"omitempty,max=50"` // Corrected: Changed from **string to *string
	SpecimenDiameter *string           `json:"specimen_diameter,omitempty"`                      // Changed to string
	SpecimenHeight   *string           `json:"specimen_height,omitempty"`                        // Changed to string
	FailureMode      *string           `json:"failure_mode,omitempty" binding:"omitempty,max=100"`
	Notes            *string           `json:"notes,omitempty"`
}

// UCSResultResponse represents the response body for UCS test result details.
type UCSResultResponse struct {
	ID               utils.BinaryUUID `json:"id"`
	SampleID         utils.BinaryUUID `json:"sample_id"`
	UCSValue         string           `json:"ucs_value"` // Changed to string
	Unit             string           `json:"unit,omitempty"`
	TestMethod       string           `json:"test_method,omitempty"`
	SpecimenDiameter string           `json:"specimen_diameter,omitempty"`
	SpecimenHeight   string           `json:"specimen_height,omitempty"`
	FailureMode      string           `json:"failure_mode,omitempty"`
	Notes            string           `json:"notes,omitempty"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at,omitempty"`
}

// ListUCSResultsResponse for listing multiple UCS results.
type ListUCSResultsResponse struct {
	UCSResults []UCSResultResponse `json:"ucs_results"`
	Total      int                 `json:"total"`
}
