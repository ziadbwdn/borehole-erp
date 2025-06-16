package models

import (
	"time"

	"boreholedata-ms/internal/utils"
)

// UserActivity represents an activity log entry for audit trails
type UserActivity struct {
	ID           utils.BinaryUUID `gorm:"primaryKey;type:binary(16)" json:"id"`    // Unique ID for the activity log entry
	UserID       utils.BinaryUUID `gorm:"type:binary(16);not null" json:"user_id"` // ID of the user who performed the action (matches User.ID)
	Username     string           `gorm:"size:50;not null" json:"username"`        // Username at the time of the action (denormalized for convenience)
	ActionType   string           `gorm:"size:50;not null" json:"action_type"`     // Type of action performed (e.g., "login", "create_project", "update_station")
	ResourceType string           `gorm:"size:50;not null" json:"resource_type"`   // Type of resource affected (e.g., "Project", "Station")
	ResourceID   *string          `gorm:"size:36" json:"resource_id,omitempty"`    // ID of the resource affected (e.g., Project ID, Station ID), optional if action is not resource-specific (e.g., login)
	Timestamp    time.Time        `gorm:"autoCreateTime" json:"timestamp"`         // When the action occurred
	IPAddress    *string          `gorm:"size:45" json:"ip_address,omitempty"`     // IP address from which the action originated, optional
	Details      *string          `gorm:"type:text" json:"details,omitempty"`      // Additional human-readable details about the action, optional
	OldValue     *string          `gorm:"type:json" json:"old_value,omitempty"`    // Optional: JSON string of the resource state *before* update
	NewValue     *string          `gorm:"type:json" json:"new_value,omitempty"`    // Optional: JSON string of the resource state *after* update
}

// ActionType constants for common user activities.
// This list can be expanded as needed to be more granular.
const (
	ActionTypeLogin                = "login"
	ActionTypeLogout               = "logout"
	ActionTypeRegisterUser         = "register_user"
	ActionTypeUpdateProfile        = "update_profile"
	ActionTypeCreateProject        = "create_project"
	ActionTypeUpdateProject        = "update_project"
	ActionTypeDeleteProject        = "delete_project"
	ActionTypeCreateStation        = "create_station"
	ActionTypeUpdateStation        = "update_station"
	ActionTypeDeleteStation        = "delete_station"
	ActionTypeUpdateDrillingStatus = "update_drilling_status" // Specific status update
	ActionTypeCreateLithology      = "create_lithology"
	ActionTypeUpdateLithology      = "update_lithology"
	ActionTypeDeleteLithology      = "delete_lithology"
	ActionTypeCreateSample         = "create_sample"
	ActionTypeUpdateSample         = "update_sample"
	ActionTypeDeleteSample         = "delete_sample"
	ActionTypeCreateLabTest        = "create_lab_test"
	ActionTypeUpdateLabTest        = "update_lab_test"
	ActionTypeDeleteLabTest        = "delete_lab_test"
	ActionTypeUpdateLabStatus      = "update_lab_status" // Specific status update for lab items
	ActionTypeCreateUCSResult      = "create_ucs_result"
	ActionTypeUpdateUCSResult      = "update_ucs_result"
	ActionTypeDeleteUCSResult      = "delete_ucs_result"
	ActionTypeGenerateReport       = "generate_report"
	ActionTypeExportData           = "export_data"
	ActionTypeFailedLogin          = "failed_login_attempt"
	ActionTypeUnauthorizedAccess   = "unauthorized_access"
	ActionTypeRoleChange           = "role_change"
)

// ResourceType constants for affected entities.
const (
	ResourceTypeUser      = "User"
	ResourceTypeProject   = "Project"
	ResourceTypeStation   = "Station"
	ResourceTypeLithology = "Lithology"
	ResourceTypeSample    = "Sample"
	ResourceTypeLabTest   = "LabTest"
	ResourceTypeUCSResult = "UCSResult"
)

// UserActivitySummary provides aggregated activity data for reporting.
type UserActivitySummary struct {
	UserID               utils.BinaryUUID `json:"userId"`
	TotalActivities      int64            `json:"totalActivities"`
	LoginCount           int64            `json:"loginCount"`
	CreateOperations     int64            `json:"createOperations"`
	UpdateOperations     int64            `json:"updateOperations"`
	DeleteOperations     int64            `json:"deleteOperations"`
	ReportGenerations    int64            `json:"reportGenerations"`
	LastActivity         time.Time        `json:"lastActivity"`
	MostAccessedResource string           `json:"mostAccessedResource"`
}
