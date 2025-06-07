package dto

import (
	"boreholedata-ms/internal/utils" // This import is correct
	"time"
)

// --- LabSample DTOs ---

// CreateLabSampleRequest represents the request body for creating a new laboratory sample.
type CreateLabSampleRequest struct {
	StationID    utils.BinaryUUID `json:"station_id" binding:"required"`
	SampleCode   string           `json:"sample_code" binding:"required,min=3,max=50"`
	DepthFrom    string           `json:"depth_from" binding:"required"` // Changed to string
	DepthTo      string           `json:"depth_to" binding:"required"`   // Changed to string
	SampleType   string           `json:"sample_type" binding:"max=50"`
	SamplingDate *time.Time       `json:"sampling_date,omitempty"` // Changed to pointer
	TestedBy     string           `json:"tested_by" binding:"max=100"`
	LabName      string           `json:"lab_name" binding:"max=100"`
	TestDate     *time.Time       `json:"test_date,omitempty"`     // Changed to pointer
}

// UpdateLabSampleRequest represents the request body for updating an existing laboratory sample.
type UpdateLabSampleRequest struct {
	StationID    *utils.BinaryUUID `json:"station_id,omitempty"`
	SampleCode   *string           `json:"sample_code,omitempty" binding:"omitempty,min=3,max=50"`
	DepthFrom    *string           `json:"depth_from,omitempty"` // Changed to string
	DepthTo      *string           `json:"depth_to,omitempty"`   // Changed to string
	SampleType   *string           `json:"sample_type,omitempty" binding:"omitempty,max=50"`
	SamplingDate *time.Time        `json:"sampling_date,omitempty"`
	TestedBy     *string           `json:"tested_by,omitempty" binding:"omitempty,max=100"`
	LabName      *string           `json:"lab_name,omitempty" binding:"omitempty,max=100"`
	TestDate     *time.Time        `json:"test_date,omitempty"`
}

// LabSampleResponse represents the response body for laboratory sample details.
type LabSampleResponse struct {
	ID           utils.BinaryUUID `json:"id"`
	StationID    utils.BinaryUUID `json:"station_id"`
	SampleCode   string           `json:"sample_code"`
	DepthFrom    string           `json:"depth_from"` // Changed to string
	DepthTo      string           `json:"depth_to"`   // Changed to string
	SampleType   string           `json:"sample_type,omitempty"`
	SamplingDate *time.Time       `json:"sampling_date,omitempty"` // Changed to pointer
	TestedBy     string           `json:"tested_by,omitempty"`
	LabName      string           `json:"lab_name,omitempty"`
	TestDate     *time.Time       `json:"test_date,omitempty"`     // Changed to pointer
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at,omitempty"`
}

// ListLabSamplesResponse for listing multiple lab samples.
type ListLabSamplesResponse struct {
	Samples []LabSampleResponse `json:"samples"`
	Total   int                 `json:"total"`
}
