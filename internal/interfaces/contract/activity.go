package contract

import (
	"context" // Import context package
	"time"

	"boreholedata-ms/internal/api/dto"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils" // Ensure this path is correct for BinaryUUID
)

// UserActivityRepository defines the interface for interacting with user activity data storage.
type UserActivityRepository interface {
	// Create logs a new user activity entry.
	Create(ctx context.Context, activity *models.UserActivity) error

	// ListUserActivities retrieves a list of user activity records based on the provided filter.
	ListUserActivities(ctx context.Context, filter dto.ActivityFilterRequest) ([]models.UserActivity, int64, error)

	// GetByUserID retrieves all activities for a specific user.
	GetByUserID(ctx context.Context, userID utils.BinaryUUID, limit, offset int) ([]*models.UserActivity, error)

	// GetByDateRange retrieves activities within a specific time period.
	GetByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*models.UserActivity, error)

	// GetByResourceType retrieves activities related to a specific resource type.
	GetByResourceType(ctx context.Context, resourceType string, limit, offset int) ([]*models.UserActivity, error) // Changed parameter name

	// GetByActionType retrieves activities of a specific action type.
	GetByActionType(ctx context.Context, actionType string, limit, offset int) ([]*models.UserActivity, error) // Changed parameter name

	// GetFailedLoginAttempts retrieves failed login attempts for security monitoring.
	GetFailedLoginAttempts(ctx context.Context, timeWindow time.Duration, limit int) ([]*models.UserActivity, error)

	// GetUserActivitySummary returns activity statistics for a user, returning a DTO.
	GetUserActivitySummary(ctx context.Context, userID utils.BinaryUUID, days int) (*dto.UserActivitySummaryResponse, error) // Now returns DTO

	// DeleteOldActivities removes activities older than specified duration.
	DeleteOldActivities(ctx context.Context, olderThan time.Duration) error
}

type UserActivityService interface {
	// LogUserActivity records a new user activity.
	LogUserActivity(ctx context.Context, userID, username, actionType, resourceType string, resourceID *string, ipAddress, details, oldValue, newValue *string) error

	// ListUserActivities fetches a paginated list of user activities based on filters.
	ListUserActivities(ctx context.Context, filter dto.ActivityFilterRequest) (*dto.UserActivityListResponse, error)

	// GetUserActivitySummary provides aggregated summary data for a specific user.
	GetUserActivitySummary(ctx context.Context, userID string, days int) (*dto.UserActivitySummaryResponse, error)

	// GetSecurityAlerts fetches specific security-related activity alerts.
	GetSecurityAlerts(ctx context.Context, timeWindow time.Duration, limit int) ([]dto.SecurityAlertResponse, error)

	// CleanOldActivities removes activities older than a specified duration.
	CleanOldActivities(ctx context.Context, olderThan time.Duration) error
}
