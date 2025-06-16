package services

import (
	"context"
	"time"

	"boreholedata-ms/internal/api/dto"
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract" // Import the contract package for the service interface
	"boreholedata-ms/internal/logger"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils" // For BinaryUUID
)

// userActivityService implements the contract.UserActivityService interface.
type userActivityService struct {
	repo   contract.UserActivityRepository // Dependency on the repository interface
	logger logger.Logger                   // Dependency on the logger
	// Add other dependencies if needed, e.g., a user service to validate UserIDs
}

// NewUserActivityService creates and returns a new instance of UserActivityService.
// It returns the interface type from the contract package.
func NewUserActivityService(repo contract.UserActivityRepository, logger logger.Logger) contract.UserActivityService {
	if repo == nil {
		panic("userActivityRepository must not be nil for UserActivityService")
	}
	if logger == nil {
		panic("logger must not be nil for UserActivityService")
	}
	return &userActivityService{
		repo:   repo,
		logger: logger,
	}
}

// LogUserActivity records a new user activity.
// It performs basic validation and constructs the model before passing it to the repository.
func (s *userActivityService) LogUserActivity(ctx context.Context, userIDStr, username, actionType, resourceType string, resourceID *string, ipAddress, details, oldValue, newValue *string) error {
	op := "userActivityService.LogUserActivity"
	s.logger.Info(ctx, "Attempting to log user activity",
		logger.Field{Key: "userID", Value: userIDStr},
		logger.Field{Key: "actionType", Value: actionType})

	// 1. Validate input parameters
	if userIDStr == "" || username == "" || actionType == "" || resourceType == "" {
		s.logger.Error(ctx, "Validation error: missing required fields for user activity", nil,
			logger.Field{Key: "userIDStr", Value: userIDStr},
			logger.Field{Key: "username", Value: username},
			logger.Field{Key: "actionType", Value: actionType},
			logger.Field{Key: "resourceType", Value: resourceType})
		return exception.NewValidationError("Missing required fields", "userID, username, actionType, and resourceType are required.")
	}

	userID, err := utils.ParseBinaryUUID(userIDStr)
	if err != nil {
		s.logger.Error(ctx, "Validation error: invalid UserID format", err, logger.Field{Key: "userIDStr", Value: userIDStr})
		return exception.NewValidationError("Invalid UserID format", err.Error())
	}

	// 2. Construct the UserActivity model
	activity := &models.UserActivity{
		UserID:       userID,
		Username:     username,
		ActionType:   actionType,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		IPAddress:    ipAddress,
		Details:      details,
		OldValue:     oldValue,
		NewValue:     newValue,
		// Timestamp will be set by the repository
		// ID will be set by the repository
	}

	// 3. Call the repository to create the activity
	if err := s.repo.Create(ctx, activity); err != nil {
		s.logger.Error(ctx, "Failed to create user activity in repository", err,
			logger.Field{Key: "userID", Value: userIDStr},
			logger.Field{Key: "actionType", Value: actionType})
		return exception.NewDatabaseError(op, err) // Wrap repository error
	}

	s.logger.Info(ctx, "User activity successfully logged",
		logger.Field{Key: "activityID", Value: activity.ID.String()},
		logger.Field{Key: "userID", Value: userIDStr},
		logger.Field{Key: "actionType", Value: actionType})
	return nil
}

// ListUserActivities fetches a paginated list of user activities based on filters.
func (s *userActivityService) ListUserActivities(ctx context.Context, filter dto.ActivityFilterRequest) (*dto.UserActivityListResponse, error) {
	op := "userActivityService.ListUserActivities"
	s.logger.Info(ctx, "Listing user activities with filter",
		logger.Field{Key: "filter", Value: filter})

	// Basic validation for pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 10 // Default page size
	} else if filter.PageSize > 100 {
		filter.PageSize = 100 // Max page size
	}

	activities, total, err := s.repo.ListUserActivities(ctx, filter)
	if err != nil {
		s.logger.Error(ctx, "Failed to list user activities from repository", err,
			logger.Field{Key: "filter", Value: filter})
		return nil, exception.NewInternalError(op, err) // Wrap repository error
	}

	// Transform models.UserActivity to dto.UserActivityResponse
	activityResponses := make([]dto.UserActivityResponse, len(activities))
	for i, activity := range activities {
		activityResponses[i] = dto.UserActivityResponse{
			ID:           activity.ID.String(),
			UserID:       activity.UserID.String(),
			Username:     activity.Username,
			ActionType:   activity.ActionType,
			ResourceType: activity.ResourceType,
			ResourceID:   activity.ResourceID,
			Timestamp:    activity.Timestamp,
			IPAddress:    activity.IPAddress,
			Details:      activity.Details,
			OldValue:     activity.OldValue,
			NewValue:     activity.NewValue,
		}
	}

	totalPages := (total + int64(filter.PageSize) - 1) / int64(filter.PageSize)

	response := &dto.UserActivityListResponse{
		Activities: activityResponses,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: int(totalPages),
	}

	s.logger.Info(ctx, "Successfully listed user activities",
		logger.Field{Key: "totalCount", Value: total},
		logger.Field{Key: "currentPage", Value: filter.Page})
	return response, nil
}

// GetUserActivitySummary provides aggregated summary data for a specific user.
func (s *userActivityService) GetUserActivitySummary(ctx context.Context, userIDStr string, days int) (*dto.UserActivitySummaryResponse, error) {
	op := "userActivityService.GetUserActivitySummary"
	s.logger.Info(ctx, "Getting user activity summary",
		logger.Field{Key: "userID", Value: userIDStr},
		logger.Field{Key: "days", Value: days})

	// Validate input
	if userIDStr == "" {
		s.logger.Error(ctx, "Validation error: userID cannot be empty for summary", nil)
		return nil, exception.NewValidationError("User ID is required", "userID cannot be empty.")
	}
	userID, err := utils.ParseBinaryUUID(userIDStr)
	if err != nil {
		s.logger.Error(ctx, "Validation error: invalid UserID format for summary", err, logger.Field{Key: "userIDStr", Value: userIDStr})
		return nil, exception.NewValidationError("Invalid UserID format", err.Error())
	}
	if days < 0 {
		days = 0 // Fetch all if days is negative, or handle as a validation error
	}

	summary, err := s.repo.GetUserActivitySummary(ctx, userID, days)
	if err != nil {
		s.logger.Error(ctx, "Failed to get user activity summary from repository", err,
			logger.Field{Key: "userID", Value: userIDStr})
		return nil, exception.NewInternalError(op, err)
	}

	s.logger.Info(ctx, "User activity summary retrieved",
		logger.Field{Key: "userID", Value: userIDStr},
		logger.Field{Key: "totalActivities", Value: summary.TotalActivities})
	return summary, nil
}

// GetSecurityAlerts fetches specific security-related activity alerts.
func (s *userActivityService) GetSecurityAlerts(ctx context.Context, timeWindow time.Duration, limit int) ([]dto.SecurityAlertResponse, error) {
	op := "userActivityService.GetSecurityAlerts"
	s.logger.Info(ctx, "Fetching security alerts",
		logger.Field{Key: "timeWindow", Value: timeWindow},
		logger.Field{Key: "limit", Value: limit})

	if timeWindow <= 0 {
		s.logger.Error(ctx, "Validation error: timeWindow must be positive for security alerts", nil)
		return nil, exception.NewValidationError("Invalid time window", "timeWindow must be a positive duration.")
	}
	if limit <= 0 {
		limit = 100 // Default limit for security alerts
	}

	activities, err := s.repo.GetFailedLoginAttempts(ctx, timeWindow, limit)
	if err != nil {
		s.logger.Error(ctx, "Failed to get failed login attempts from repository", err)
		return nil, exception.NewInternalError(op, err)
	}

	alerts := make([]dto.SecurityAlertResponse, 0, len(activities))
	for _, activity := range activities {
		// Example of mapping a failed login activity to a SecurityAlertResponse
		if activity.ActionType == models.ActionTypeFailedLogin {
			alert := dto.SecurityAlertResponse{
				AlertID:     utils.NewBinaryUUID().String(), // Generate new ID for the alert DTO
				UserID:      activity.UserID.String(),
				Username:    activity.Username,
				ActionType:  activity.ActionType,
				Timestamp:   activity.Timestamp,
				IPAddress:   activity.IPAddress,
				Description: "Repeated failed login attempt detected.",
				Severity:    "High",
				RelatedActivity: &dto.UserActivityResponse{ // Optionally include the full activity
					ID:           activity.ID.String(),
					UserID:       activity.UserID.String(),
					Username:     activity.Username,
					ActionType:   activity.ActionType,
					ResourceType: activity.ResourceType,
					ResourceID:   activity.ResourceID,
					Timestamp:    activity.Timestamp,
					IPAddress:    activity.IPAddress,
					Details:      activity.Details,
					OldValue:     activity.OldValue,
					NewValue:     activity.NewValue,
				},
			}
			alerts = append(alerts, alert)
		}
	}

	s.logger.Info(ctx, "Security alerts fetched successfully",
		logger.Field{Key: "alertCount", Value: len(alerts)})
	return alerts, nil
}

// CleanOldActivities removes activities older than a specified duration.
func (s *userActivityService) CleanOldActivities(ctx context.Context, olderThan time.Duration) error {
	op := "userActivityService.CleanOldActivities"
	s.logger.Info(ctx, "Initiating cleanup of old user activities",
		logger.Field{Key: "olderThan", Value: olderThan.String()})

	if olderThan <= 0 {
		s.logger.Error(ctx, "Validation error: olderThan duration must be positive for cleanup", nil)
		return exception.NewValidationError("Invalid duration", "olderThan duration must be a positive value.")
	}

	if err := s.repo.DeleteOldActivities(ctx, olderThan); err != nil {
		s.logger.Error(ctx, "Failed to delete old activities from repository", err,
			logger.Field{Key: "olderThan", Value: olderThan.String()})
		return exception.NewInternalError(op, err)
	}

	s.logger.Info(ctx, "Old user activities cleaned up successfully",
		logger.Field{Key: "olderThan", Value: olderThan.String()})
	return nil
}
