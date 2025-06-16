package dto

import (
	"time"
)

// LogActivityRequest represents the request body for logging a new user activity.
type LogActivityRequest struct {
	Username     string  `json:"username" binding:"required"`      // Username of the user performing the action
	ActionType   string  `json:"action_type" binding:"required"`   // Type of action (e.g., "login", "create_project")
	ResourceType string  `json:"resource_type" binding:"required"` // Type of resource (e.g., "Project", "Station")
	ResourceID   *string `json:"resource_id,omitempty"`            // ID of the resource affected, optional
	IPAddress    *string `json:"ip_address,omitempty"`             // IP address from which the action originated, optional
	Details      *string `json:"details,omitempty"`                // Additional human-readable details, optional
	OldValue     *string `json:"old_value,omitempty"`              // JSON string of resource state before update, optional
	NewValue     *string `json:"new_value,omitempty"`              // JSON string of resource state after update, optional
}

// UserActivityResponse represents the response structure for a single user activity entry.
// It mirrors essential fields from the UserActivity model for API consumption.
type UserActivityResponse struct {
	ID           string    `json:"id"`                   // UUID as string
	UserID       string    `json:"userId"`               // User's UUID as string
	Username     string    `json:"username"`             // Denormalized username for display
	ActionType   string    `json:"actionType"`           // e.g., "login", "create_project"
	ResourceType string    `json:"resourceType"`         // e.g., "Project", "User"
	ResourceID   *string   `json:"resourceId,omitempty"` // ID of the affected resource, optional
	Timestamp    time.Time `json:"timestamp"`            // When the action occurred
	IPAddress    *string   `json:"ipAddress,omitempty"`  // IP address of the requester, optional
	Details      *string   `json:"details,omitempty"`    // Additional human-readable details, optional
	OldValue     *string   `json:"oldValue,omitempty"`   // JSON string of resource state before change, optional
	NewValue     *string   `json:"newValue,omitempty"`   // JSON string of resource state after change, optional
}

// UserActivityListResponse represents a paginated list of user activities.
type UserActivityListResponse struct {
	Activities []UserActivityResponse `json:"activities"` // Slice of user activity entries
	Total      int64                  `json:"total"`      // Total number of activities matching criteria
	Page       int                    `json:"page"`       // Current page number
	PageSize   int                    `json:"pageSize"`   // Number of items per page
	TotalPages int                    `json:"totalPages"` // Total number of pages
}

// UserActivitySummaryResponse provides aggregated activity data for reporting.
type UserActivitySummaryResponse struct {
	UserID               string    `json:"userId"`               // User ID as string
	TotalActivities      int64     `json:"totalActivities"`      // Total number of activities for the user
	LoginCount           int64     `json:"loginCount"`           // Count of login activities
	CreateOperations     int64     `json:"createOperations"`     // Count of 'create' actions across resources
	UpdateOperations     int64     `json:"updateOperations"`     // Count of 'update' actions
	DeleteOperations     int64     `json:"deleteOperations"`     // Count of 'delete' actions
	ReportGenerations    int64     `json:"reportGenerations"`    // Count of report generation actions
	LastActivity         time.Time `json:"lastActivity"`         // Timestamp of the last activity
	MostAccessedResource string    `json:"mostAccessedResource"` // The resource type most frequently accessed
	// You can add more aggregated fields here as needed, e.g.:
	// ActivitiesByActionType   map[string]int64 `json:"activitiesByActionType"`
	// ActivitiesByResourceType map[string]int64 `json:"activitiesByResourceType"`
}

// ActivityFilterRequest represents the request parameters for filtering user activities.
type ActivityFilterRequest struct {
	UserID       *string    `json:"userId,omitempty" form:"userId"`                                               // Filter by specific user ID
	ActionType   *string    `json:"actionType,omitempty" form:"actionType"`                                       // Filter by type of action (e.g., "login")
	ResourceType *string    `json:"resourceType,omitempty" form:"resourceType"`                                   // Filter by type of resource (e.g., "Project")
	ResourceID   *string    `json:"resourceId,omitempty" form:"resourceId"`                                       // Filter by specific resource ID
	IPAddress    *string    `json:"ipAddress,omitempty" form:"ipAddress"`                                         // Filter by IP address
	StartDate    *time.Time `json:"startDate,omitempty" form:"startDate" time_format:"2006-01-02T15:04:05Z07:00"` // Start of time range (ISO 8601)
	EndDate      *time.Time `json:"endDate,omitempty" form:"endDate" time_format:"2006-01-02T15:04:05Z07:00"`     // End of time range (ISO 8601)
	SearchTerm   *string    `json:"searchTerm,omitempty" form:"searchTerm"`                                       // Generic search across fields like details or username

	Page     int `json:"page" form:"page,default=1"`          // Page number for pagination, defaults to 1
	PageSize int `json:"pageSize" form:"pageSize,default=10"` // Number of items per page, defaults to 10
}

// SecurityAlertResponse represents a specific type of user activity that indicates a potential security concern.
type SecurityAlertResponse struct {
	AlertID         string                `json:"alertId"`                   // Unique ID for the alert itself
	UserID          string                `json:"userId"`                    // ID of the user involved in the alert
	Username        string                `json:"username"`                  // Username of the user involved
	ActionType      string                `json:"actionType"`                // e.g., "failed_login_attempt", "unauthorized_access"
	Timestamp       time.Time             `json:"timestamp"`                 // Time when the alertable event occurred
	IPAddress       *string               `json:"ipAddress,omitempty"`       // IP address related to the alert
	Description     string                `json:"description"`               // A concise description of the security event
	Severity        string                `json:"severity"`                  // e.g., "High", "Medium", "Low"
	RelatedActivity *UserActivityResponse `json:"relatedActivity,omitempty"` // Optional: Link to the full activity record
}
