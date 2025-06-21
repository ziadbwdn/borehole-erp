package repository

import (
	"context"
	"strings"
	"time"
	"fmt"
	"boreholedata-ms/internal/api/dto"
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/logger"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils" // For BinaryUUID

	"gorm.io/gorm"
)

// GormUserActivityRepository implements UserActivityRepository for a GORM database.
type GormUserActivityRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewGormUserActivityRepository creates a new GormUserActivityRepository instance.
// It requires a *gorm.DB instance and a logger.
func NewGormUserActivityRepository(db *gorm.DB, logger logger.Logger) *GormUserActivityRepository {
	if db == nil {
		panic("gorm.DB instance must not be nil for GormUserActivityRepository")
	}
	if logger == nil {
		panic("logger must not be nil for GormUserActivityRepository")
	}
	return &GormUserActivityRepository{
		db:     db,
		logger: logger,
	}
}

// Create saves a new user activity record to the database.
func (r *GormUserActivityRepository) Create(ctx context.Context, activity *models.UserActivity) error {
	r.logger.Info(ctx, "GORMRepo: Attempting to save user activity",
		logger.Field{Key: "proposedActivityID", Value: activity.ID.String()}, // Log the ID before creation
		logger.Field{Key: "userID", Value: activity.UserID.String()},
		logger.Field{Key: "actionType", Value: activity.ActionType},
		logger.Field{Key: "username", Value: activity.Username},
		logger.Field{Key: "resourceID", Value: activity.ResourceID}, // Show *string value
		logger.Field{Key: "ipAddress", Value: activity.IPAddress},   // Show *string value
	)

	result := r.db.WithContext(ctx).Create(activity) // Store the result

	if result.Error != nil {
		r.logger.Error(ctx, "GORMRepo: Failed to create user activity in DB", result.Error,
			logger.Field{Key: "userID", Value: activity.UserID.String()},
			logger.Field{Key: "actionType", Value: activity.ActionType},
			logger.Field{Key: "rowsAffected", Value: result.RowsAffected}, // Log RowsAffected even on error
		)
		return exception.NewDatabaseError("CreateUserActivity", result.Error)
	}

	if result.RowsAffected == 0 {
		// This is a very suspicious case for a Create operation
		r.logger.Warn(ctx, "GORMRepo: Create user activity reported 0 rows affected with no error",
			logger.Field{Key: "activityID", Value: activity.ID.String()},
			logger.Field{Key: "userID", Value: activity.UserID.String()},
			logger.Field{Key: "actionType", Value: activity.ActionType},
		)
		// Consider this an error as well, or at least investigate why it happened
		return exception.NewDatabaseError("CreateUserActivity", fmt.Errorf("0 rows affected on create"))
	}

	r.logger.Info(ctx, "GORMRepo: User activity created in DB successfully",
		logger.Field{Key: "activityID", Value: activity.ID.String()},
		logger.Field{Key: "userID", Value: activity.UserID.String()},
		logger.Field{Key: "actionType", Value: activity.ActionType},
		logger.Field{Key: "rowsAffected", Value: result.RowsAffected},
	)
	return nil
}

// ListUserActivities retrieves a list of user activity records based on the provided filter.
func (r *GormUserActivityRepository) ListUserActivities(ctx context.Context, filter dto.ActivityFilterRequest) ([]models.UserActivity, int64, error) {
	var activities []models.UserActivity
	query := r.db.WithContext(ctx).Model(&models.UserActivity{})

	// Apply filters
	if filter.UserID != nil {
		// Ensure filter.UserID is correctly converted to the database's UUID type if needed
		parsedUserID, err := utils.ParseBinaryUUID(*filter.UserID) // Assuming filter.UserID is string form
		if err != nil {
			r.logger.Error(ctx, "Invalid UserID filter format", err, logger.Field{Key: "rawUserID", Value: *filter.UserID})
			return nil, 0, exception.NewValidationError("Invalid UserID format for filter")
		}
		query = query.Where("user_id = ?", parsedUserID)
	}
	if filter.ActionType != nil {
		query = query.Where("action_type = ?", *filter.ActionType)
	}
	if filter.ResourceType != nil {
		query = query.Where("resource_type = ?", *filter.ResourceType)
	}
	if filter.ResourceID != nil {
		query = query.Where("resource_id = ?", *filter.ResourceID)
	}
	if filter.IPAddress != nil {
		query = query.Where("ip_address LIKE ?", "%"+strings.ToLower(*filter.IPAddress)+"%")
	}
	if filter.StartDate != nil {
		query = query.Where("timestamp >= ?", filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("timestamp <= ?", filter.EndDate)
	}
	if filter.SearchTerm != nil {
		searchTermLower := "%" + strings.ToLower(*filter.SearchTerm) + "%"
		query = query.Where("LOWER(details) LIKE ? OR LOWER(username) LIKE ?", searchTermLower, searchTermLower)
	}

	var totalCount int64
	// Get total count before applying limit and offset
	if err := query.Count(&totalCount).Error; err != nil {
		r.logger.Error(ctx, "Failed to count user activities", err)
		return nil, 0, exception.NewDatabaseError("ListUserActivities.Count", err)
	}

	// Apply sorting (most recent first)
	query = query.Order("timestamp DESC")

	// Apply pagination
	if filter.PageSize > 0 && filter.Page > 0 {
		query = query.Limit(int(filter.PageSize)).Offset(int((filter.Page - 1) * filter.PageSize))
	}

	if err := query.Find(&activities).Error; err != nil {
		r.logger.Error(ctx, "Failed to list user activities from DB", err,
			logger.Field{Key: "filter", Value: filter})
		return nil, 0, exception.NewDatabaseError("ListUserActivities.Find", err)
	}

	r.logger.Info(ctx, "Activities listed from DB successfully",
		logger.Field{Key: "totalFiltered", Value: totalCount},
		logger.Field{Key: "page", Value: filter.Page},
		logger.Field{Key: "pageSize", Value: filter.PageSize},
		logger.Field{Key: "returnedCount", Value: len(activities)})

	return activities, totalCount, nil
}

// GetByUserID retrieves all activities for a specific user.
func (r *GormUserActivityRepository) GetByUserID(ctx context.Context, userID utils.BinaryUUID, limit, offset int) ([]*models.UserActivity, error) {
	var activities []*models.UserActivity
	query := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&activities).Error; err != nil {
		r.logger.Error(ctx, "Failed to get activities by UserID from DB", err, logger.Field{Key: "userID", Value: userID.String()})
		return nil, exception.NewDatabaseError("GetByUserID", err)
	}
	r.logger.Info(ctx, "Activities retrieved by UserID from DB", logger.Field{Key: "userID", Value: userID.String()}, logger.Field{Key: "count", Value: len(activities)})
	return activities, nil
}

// GetByDateRange retrieves activities within a specific time period.
func (r *GormUserActivityRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*models.UserActivity, error) {
	var activities []*models.UserActivity
	query := r.db.WithContext(ctx).
		Where("timestamp BETWEEN ? AND ?", startDate, endDate).
		Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&activities).Error; err != nil {
		r.logger.Error(ctx, "Failed to get activities by date range from DB", err,
			logger.Field{Key: "startDate", Value: startDate.Format(time.RFC3339)},
			logger.Field{Key: "endDate", Value: endDate.Format(time.RFC3339)})
		return nil, exception.NewDatabaseError("GetByDateRange", err)
	}
	r.logger.Info(ctx, "Activities retrieved by date range from DB", logger.Field{Key: "count", Value: len(activities)})
	return activities, nil
}

// GetByResourceType retrieves activities related to a specific resource type.
func (r *GormUserActivityRepository) GetByResourceType(ctx context.Context, resourceType string, limit, offset int) ([]*models.UserActivity, error) {
	var activities []*models.UserActivity
	query := r.db.WithContext(ctx).
		Where("LOWER(resource_type) = LOWER(?)", resourceType).
		Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&activities).Error; err != nil {
		r.logger.Error(ctx, "Failed to get activities by resource type from DB", err, logger.Field{Key: "resourceType", Value: resourceType})
		return nil, exception.NewDatabaseError("GetByResourceType", err)
	}
	r.logger.Info(ctx, "Activities retrieved by resource type from DB", logger.Field{Key: "resourceType", Value: resourceType}, logger.Field{Key: "count", Value: len(activities)})
	return activities, nil
}

// GetByActionType retrieves activities of a specific action type.
func (r *GormUserActivityRepository) GetByActionType(ctx context.Context, actionType string, limit, offset int) ([]*models.UserActivity, error) {
	var activities []*models.UserActivity
	query := r.db.WithContext(ctx).
		Where("LOWER(action_type) = LOWER(?)", actionType).
		Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&activities).Error; err != nil {
		r.logger.Error(ctx, "Failed to get activities by action type from DB", err, logger.Field{Key: "actionType", Value: actionType})
		return nil, exception.NewDatabaseError("GetByActionType", err)
	}
	r.logger.Info(ctx, "Activities retrieved by action type from DB", logger.Field{Key: "actionType", Value: actionType}, logger.Field{Key: "count", Value: len(activities)})
	return activities, nil
}

// GetFailedLoginAttempts retrieves failed login attempts within a specified time window.
func (r *GormUserActivityRepository) GetFailedLoginAttempts(ctx context.Context, timeWindow time.Duration, limit int) ([]*models.UserActivity, error) {
	var failedLogins []*models.UserActivity
	threshold := time.Now().Add(-timeWindow)

	query := r.db.WithContext(ctx).
		Where("action_type = ? AND timestamp >= ?", models.ActionTypeFailedLogin, threshold).
		Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&failedLogins).Error; err != nil {
		r.logger.Error(ctx, "Failed to get failed login attempts from DB", err, logger.Field{Key: "timeWindow", Value: timeWindow.String()})
		return nil, exception.NewDatabaseError("GetFailedLoginAttempts", err)
	}
	r.logger.Info(ctx, "Failed login attempts retrieved from DB",
		logger.Field{Key: "timeWindow", Value: timeWindow.String()},
		logger.Field{Key: "returnedCount", Value: len(failedLogins)})
	return failedLogins, nil
}

// GetUserActivitySummary returns activity statistics for a user for a given number of days.
func (r *GormUserActivityRepository) GetUserActivitySummary(ctx context.Context, userID utils.BinaryUUID, days int) (*dto.UserActivitySummaryResponse, error) {
	// This method requires aggregation which is often complex in raw GORM,
	// especially for "MostAccessedResource" which needs group by and order by count.
	// For simplicity, we'll fetch all relevant activities and process in-memory,
	// or perform multiple targeted queries.
	// For production, you might consider direct SQL queries or more advanced GORM features/plugins.

	var activities []models.UserActivity
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if days > 0 {
		cutoff := time.Now().Add(time.Duration(-days) * 24 * time.Hour)
		query = query.Where("timestamp >= ?", cutoff)
	}

	if err := query.Find(&activities).Error; err != nil {
		r.logger.Error(ctx, "Failed to get user activities for summary from DB", err,
			logger.Field{Key: "userID", Value: userID.String()},
			logger.Field{Key: "days", Value: days})
		return nil, exception.NewDatabaseError("GetUserActivitySummary", err)
	}

	summary := &dto.UserActivitySummaryResponse{
		UserID:               userID.String(),
		TotalActivities:      int64(len(activities)),
		LoginCount:           0,
		CreateOperations:     0,
		UpdateOperations:     0,
		DeleteOperations:     0,
		ReportGenerations:    0,
		ExportGenerations:	  0,
		LastActivity:         time.Time{}, // Zero value
		MostAccessedResource: "",
	}

	var lastActivity time.Time
	resourceCounts := make(map[string]int64)

	for _, activity := range activities {
		if lastActivity.IsZero() || activity.Timestamp.After(lastActivity) {
			lastActivity = activity.Timestamp
		}

		switch activity.ActionType {
		case models.ActionTypeLogin:
			summary.LoginCount++
		case models.ActionTypeCreateProject, models.ActionTypeCreateStation, models.ActionTypeCreateLithology,
			models.ActionTypeCreateSample, models.ActionTypeCreateLabTest, models.ActionTypeCreateUCSResult:
			summary.CreateOperations++
		case models.ActionTypeUpdateProject, models.ActionTypeUpdateStation, models.ActionTypeUpdateDrillingStatus,
			models.ActionTypeUpdateLithology, models.ActionTypeUpdateSample, models.ActionTypeUpdateLabTest,
			models.ActionTypeUpdateLabStatus, models.ActionTypeUpdateUCSResult, models.ActionTypeUpdateProfile:
			summary.UpdateOperations++
		case models.ActionTypeDeleteProject, models.ActionTypeDeleteStation, models.ActionTypeDeleteLithology,
			models.ActionTypeDeleteSample, models.ActionTypeDeleteLabTest, models.ActionTypeDeleteUCSResult:
			summary.DeleteOperations++
		case models.ActionTypeGenerateReport:
			summary.ReportGenerations++
		case models.ActionTypeExportData:
			summary.ExportGenerations++
		}

		if activity.ResourceType != "" {
			resourceCounts[activity.ResourceType]++
		}
	}

	summary.LastActivity = lastActivity

	maxCount := int64(0)
	mostAccessedResource := ""
	for resource, count := range resourceCounts {
		if count > maxCount {
			maxCount = count
			mostAccessedResource = resource
		}
	}
	summary.MostAccessedResource = mostAccessedResource

	r.logger.Info(ctx, "User activity summary generated from DB",
		logger.Field{Key: "userID", Value: userID.String()},
		logger.Field{Key: "days", Value: days},
		logger.Field{Key: "totalActivities", Value: summary.TotalActivities})
	return summary, nil
}

// DeleteOldActivities removes activities older than the specified duration.
func (r *GormUserActivityRepository) DeleteOldActivities(ctx context.Context, olderThan time.Duration) error {
	threshold := time.Now().Add(-olderThan)
	result := r.db.WithContext(ctx).Where("timestamp < ?", threshold).Delete(&models.UserActivity{})
	if result.Error != nil {
		r.logger.Error(ctx, "Failed to delete old activities from DB", result.Error,
			logger.Field{Key: "olderThan", Value: olderThan.String()})
		return exception.NewDatabaseError("DeleteOldActivities", result.Error)
	}

	r.logger.Info(ctx, "Old activities deleted from DB",
		logger.Field{Key: "olderThan", Value: olderThan.String()},
		logger.Field{Key: "deletedCount", Value: result.RowsAffected})
	return nil
}
