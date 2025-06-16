package dto

import (
	"time"

	"boreholedata-ms/internal/utils" // Import utils for BinaryUUID
)

// CreateProjectRequest represents the request body for creating a new project.
type CreateProjectRequest struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartDate   time.Time `json:"startDate,omitzero"` // Changed to pointer and removed not null
	EndDate     time.Time `json:"endDate,omitzero"`   // Changed to pointer and removed not null
}

// UpdateProjectRequest represents the request body for updating an existing project.
// Fields are pointers to allow partial updates (nil means no change).
type UpdateProjectRequest struct {
	Name        *string    `json:"name,omitempty" binding:"omitempty,min=3"`
	Description *string    `json:"description,omitempty"`
	Location    *string    `json:"location,omitempty"`
	Status      *string    `json:"status,omitempty"`    // Assuming status can be updated
	StartDate   *time.Time `json:"startDate,omitempty"` // Explicitly use pointer
	EndDate     *time.Time `json:"endDate,omitempty"`   // Explicitly use pointer
}

// ProjectResponse represents the response body for project details.
type ProjectResponse struct {
	ID          utils.BinaryUUID `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"` // Include description in response
	Location    string           `json:"location"`
	StartDate   *time.Time       `json:"startDate,omitempty"` // Changed to pointer
	EndDate     *time.Time       `json:"endDate,omitempty"`   // Changed to pointer
	Status      string           `json:"status"`
	CreatedBy   utils.BinaryUUID `json:"createdBy"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt,omitempty"`
}

// ListProjectsResponse (optional, for listing multiple projects)
type ListProjectsResponse struct {
	Projects []ProjectResponse `json:"projects"`
	Total    int               `json:"total"`
}
