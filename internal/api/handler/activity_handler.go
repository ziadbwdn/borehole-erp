package handler

import (
	"net/http"
	"strconv"
	"time"

	"boreholedata-ms/internal/api/dto"
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/internal/logger"
	"boreholedata-ms/pkg/gin_helpers"
	"boreholedata-ms/pkg/http_response"

	"github.com/gin-gonic/gin"
)

// UserActivityHandler defines the HTTP handlers for user activity operations.
type UserActivityHandler struct {
	service contract.UserActivityService
	logger  logger.Logger
}

// NewUserActivityHandler creates and returns a new instance of UserActivityHandler.
func NewUserActivityHandler(service contract.UserActivityService, logger logger.Logger) *UserActivityHandler {
	if service == nil {
		panic("userActivityService must not be nil for UserActivityHandler")
	}
	if logger == nil {
		panic("logger must not be nil for UserActivityHandler")
	}
	return &UserActivityHandler{
		service: service,
		logger:  logger,
	}
}

// LogActivity godoc
// @Router /activities [post]
func (h *UserActivityHandler) LogActivity(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.LogActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx, "Failed to bind log activity request", err) // Correct: err is passed
		http_response.HandleAppError(c, exception.NewValidationError("Invalid request body", err.Error()))
		return
	}

	userID, appErr := gin_helpers.GetUserIDFromContext(c)
	if appErr != nil {
		h.logger.Error(ctx, "Failed to get userID from context for logging activity", appErr) // Correct: appErr is passed as error
		http_response.HandleAppError(c, appErr)
		return
	}

	username := req.Username
	if username == "" {
		usernameAny, exists := c.Get("username")
		if exists {
			if u, ok := usernameAny.(string); ok {
				username = u
			}
		}
		if username == "" {
			h.logger.Info(ctx, "Username not found in request or context, using 'unknown' placeholder") // Use Info or specific error
			username = "unknown"
		}
	}

	err := h.service.LogUserActivity(
		ctx,
		userID.String(),
		username,
		req.ActionType,
		req.ResourceType,
		req.ResourceID,
		req.IPAddress,
		req.Details,
		req.OldValue,
		req.NewValue,
	)

	if err != nil {
		h.logger.Error(ctx, "Failed to log user activity via service", err) // Correct: err is passed
		if appErr, ok := err.(*exception.AppError); ok {
			http_response.HandleAppError(c, appErr)
		} else {
			http_response.HandleAppError(c, exception.NewInternalError("Failed to log activity", err))
		}
		return
	}

	http_response.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Activity logged successfully"})
}

// ListActivities godoc
// @Router /activities [get]
func (h *UserActivityHandler) ListActivities(c *gin.Context) {
	ctx := c.Request.Context()

	var filter dto.ActivityFilterRequest

	if err := c.ShouldBindQuery(&filter); err != nil {
		h.logger.Error(ctx, "Failed to bind query parameters for listing activities", err) // Correct: err is passed
		http_response.HandleAppError(c, exception.NewValidationError("Invalid query parameters", err.Error()))
		return
	}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			h.logger.Error(ctx, "Failed to parse start_date", err, logger.Field{Key: "startDate", Value: startDateStr}) // Correct: err passed, fields used
			http_response.HandleAppError(c, exception.NewValidationError("Invalid start_date format", "Expected YYYY-MM-DDTHH:MM:SSZ"))
			return
		}
		filter.StartDate = &parsedTime
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			h.logger.Error(ctx, "Failed to parse end_date", err, logger.Field{Key: "endDate", Value: endDateStr}) // Correct: err passed, fields used
			http_response.HandleAppError(c, exception.NewValidationError("Invalid end_date format", "Expected YYYY-MM-DDTHH:MM:SSZ"))
			return
		}
		filter.EndDate = &parsedTime
	}

	listResponse, err := h.service.ListUserActivities(ctx, filter)
	if err != nil {
		h.logger.Error(ctx, "Failed to list user activities via service", err, logger.Field{Key: "filter", Value: filter}) // Correct: err passed, fields used
		if appErr, ok := err.(*exception.AppError); ok {
			http_response.HandleAppError(c, appErr)
		} else {
			http_response.HandleAppError(c, exception.NewInternalError("Failed to retrieve activities", err))
		}
		return
	}

	http_response.RespondWithSuccess(c, http.StatusOK, listResponse)
}

// GetActivitySummary godoc
// @Router /activities/summary/{userID} [get]
func (h *UserActivityHandler) GetActivitySummary(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr := c.Param("userID")
	if userIDStr == "" {
		h.logger.Error(ctx, "Validation error: UserID path parameter missing for summary", nil) // Correct: nil for error, no fields
		http_response.HandleAppError(c, exception.NewValidationError("User ID is required", "User ID cannot be empty."))
		return
	}

	daysStr := c.DefaultQuery("days", "0")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		h.logger.Error(ctx, "Failed to parse 'days' query parameter for summary", err, logger.Field{Key: "daysStr", Value: daysStr}) // Correct: err passed, fields used
		http_response.HandleAppError(c, exception.NewValidationError("Invalid 'days' parameter", "Days must be an integer."))
		return
	}

	summary, err := h.service.GetUserActivitySummary(ctx, userIDStr, days)
	if err != nil {
		h.logger.Error(ctx, "Failed to get user activity summary via service", err, logger.Field{Key: "userID", Value: userIDStr}) // Correct: err passed, fields used
		if appErr, ok := err.(*exception.AppError); ok {
			http_response.HandleAppError(c, appErr)
		} else {
			http_response.HandleAppError(c, exception.NewInternalError("Failed to retrieve activity summary", err))
		}
		return
	}

	http_response.RespondWithSuccess(c, http.StatusOK, summary)
}

// GetSecurityAlerts godoc
// @Router /activities/alerts [get]
func (h *UserActivityHandler) GetSecurityAlerts(c *gin.Context) {
	ctx := c.Request.Context()

	timeWindowStr := c.DefaultQuery("time_window", "24h")
	timeWindow, err := time.ParseDuration(timeWindowStr)
	if err != nil {
		h.logger.Error(ctx, "Failed to parse 'time_window' query parameter for alerts", err, logger.Field{Key: "timeWindowStr", Value: timeWindowStr}) // Correct: err passed, fields used
		http_response.HandleAppError(c, exception.NewValidationError("Invalid time_window format", "Expected a duration (e.g., 1h, 24h)."))
		return
	}

	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.logger.Error(ctx, "Failed to parse 'limit' query parameter for alerts", err, logger.Field{Key: "limitStr", Value: limitStr}) // Correct: err passed, fields used
		http_response.HandleAppError(c, exception.NewValidationError("Invalid 'limit' parameter", "Limit must be an integer."))
		return
	}

	alerts, err := h.service.GetSecurityAlerts(ctx, timeWindow, limit)
	if err != nil {
		h.logger.Error(ctx, "Failed to get security alerts via service", err) // Correct: err passed, no extra fields
		if appErr, ok := err.(*exception.AppError); ok {
			http_response.HandleAppError(c, appErr)
		} else {
			http_response.HandleAppError(c, exception.NewInternalError("Failed to retrieve security alerts", err))
		}
		return
	}

	http_response.RespondWithSuccess(c, http.StatusOK, alerts)
}

// CleanOldActivities godoc
// @Security BearerAuth // Assuming this endpoint requires authentication/authorization
func (h *UserActivityHandler) CleanOldActivities(c *gin.Context) {
	ctx := c.Request.Context()

	olderThanStr := c.Query("older_than")
	if olderThanStr == "" {
		h.logger.Error(ctx, "Validation error: 'older_than' query parameter missing for cleanup", nil) // Correct: nil for error, no fields
		http_response.HandleAppError(c, exception.NewValidationError("Missing parameter", "'older_than' duration is required."))
		return
	}

	olderThan, err := time.ParseDuration(olderThanStr)
	if err != nil {
		h.logger.Error(ctx, "Failed to parse 'older_than' query parameter for cleanup", err, logger.Field{Key: "olderThanStr", Value: olderThanStr}) // Correct: err passed, fields used
		http_response.HandleAppError(c, exception.NewValidationError("Invalid 'older_than' format", "Expected a duration (e.g., 720h, 30m)."))
		return
	}

	err = h.service.CleanOldActivities(ctx, olderThan)
	if err != nil {
		h.logger.Error(ctx, "Failed to clean old activities via service", err) // Correct: err passed, no extra fields
		if appErr, ok := err.(*exception.AppError); ok {
			http_response.HandleAppError(c, appErr)
		} else {
			http_response.HandleAppError(c, exception.NewInternalError("Failed to clean activities", err))
		}
		return
	}

	http_response.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Old activities cleaned successfully"})
}
