package services

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/internal/logger" // <-- Added
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"context"
	"encoding/json" // <-- Added
	"fmt"
	"time"
)

// LithologyServiceImpl implements the contract.LithologyService interface.
// REFACTORED: Added activityService and logger fields.
type LithologyServiceImpl struct {
	lithologyRepo   contract.LithologyRepository
	stationRepo     contract.StationRepository
	activityService contract.UserActivityService
	logger          logger.Logger
}

// NewLithologyService creates and returns a new instance of LithologyServiceImpl.
// REFACTORED: Now accepts UserActivityService and Logger.
func NewLithologyService(
	lithologyRepo contract.LithologyRepository,
	stationRepo contract.StationRepository,
	activityService contract.UserActivityService,
	logger logger.Logger,
) contract.LithologyService {
	// Panic checks
	if lithologyRepo == nil { panic("lithologyRepo must not be nil") }
	if stationRepo == nil { panic("stationRepo must not be nil") }
	if activityService == nil { panic("activityService must not be nil") }
	if logger == nil { panic("logger must not be nil") }

	return &LithologyServiceImpl{
		lithologyRepo:   lithologyRepo,
		stationRepo:     stationRepo,
		activityService: activityService,
		logger:          logger,
	}
}

// CreateLog handles the business logic for creating a new lithology log.
// REFACTORED: Added user activity logging on success.
func (s *LithologyServiceImpl) CreateLog(ctx context.Context, log *models.LithologyLog, createdBy utils.BinaryUUID, userRole models.UserRole, logCtx models.ActivityLogContext) (*models.LithologyLog, *exception.AppError) {
	if userRole != models.RoleGeologist {
		return nil, exception.NewPermissionError("User not authorized to create lithology logs")
	}
	_, appErr := s.stationRepo.GetByID(ctx, log.StationID)
	if appErr != nil {
		return nil, exception.NewValidationError(fmt.Sprintf("Station with ID '%s' not found", log.StationID.String()))
	}
	log.ID = utils.NewBinaryUUID()
	log.CreatedAt = time.Now()
	log.UpdatedAt = time.Now()

	appErr = s.lithologyRepo.CreateLog(ctx, log)
	if appErr != nil {
		return nil, appErr
	}

	// --- LOG USER ACTIVITY ---
	newValueJSON, _ := json.Marshal(log)
	newValueStr := string(newValueJSON)
	logIDStr := log.ID.String()
	details := fmt.Sprintf("New Lithology Log created for Station ID %s.", log.StationID.String())
	ipAddr := logCtx.IPAddress

	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, models.ActionTypeCreateLithology, models.ResourceTypeLithology, &logIDStr, &ipAddr, &details, nil, &newValueStr)
	
	if logErr != nil {
		s.logger.Error(ctx, "Failed to log CreateLog activity", logErr, logger.Field{Key: "logID", Value: log.ID.String()})
	}
	
	return log, nil
}
// GetLogByID retrieves a lithology log by its ID.
func (s *LithologyServiceImpl) GetLogByID(
	ctx context.Context,
	id utils.BinaryUUID,
) (*models.LithologyLog, *exception.AppError) {
	log, appErr := s.lithologyRepo.GetLogByID(ctx, id)
	if appErr != nil {
		return nil, appErr
	}
	// Optional logging can be added here if needed, following the same pattern.
	return log, nil
}

// UpdateLog handles the business logic for updating an existing lithology log.
func (s *LithologyServiceImpl) UpdateLog(ctx context.Context, log *models.LithologyLog, userRole models.UserRole, logCtx models.ActivityLogContext) *exception.AppError {
	if userRole != models.RoleGeologist {
		return exception.NewPermissionError("User not authorized to update lithology logs")
	}
	existingLog, appErr := s.lithologyRepo.GetLogByID(ctx, log.ID)
	if appErr != nil {
		return appErr
	}
	oldValueJSON, err := json.Marshal(existingLog)
	if err != nil {
		s.logger.Warn(ctx, "Failed to marshal old lithology log value for logging", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "logID", Value: log.ID.String()})
	}
	oldValueStr := string(oldValueJSON)
	// --- END CAPTURE ---

	// Apply updates from request
	if log.StationID != (utils.BinaryUUID{}) {
		_, stationErr := s.stationRepo.GetByID(ctx, log.StationID)
		if stationErr != nil {
			return exception.NewValidationError(fmt.Sprintf("New Station with ID '%s' not found or inaccessible", log.StationID.String()))
		}
		existingLog.StationID = log.StationID
	}
	// ... [applying all other fields from 'log' to 'existingLog' as in the original code] ...
	if log.DepthFrom.Internal.Value != "" { existingLog.DepthFrom.Internal.Value = log.DepthFrom.Internal.Value }
	if log.DepthTo.Internal.Value != "" { existingLog.DepthTo.Internal.Value = log.DepthTo.Internal.Value }
	if log.LithologyType != "" { existingLog.LithologyType = log.LithologyType }
	if log.RockColor != "" { existingLog.RockColor = log.RockColor }
	if log.GrainSize != "" { existingLog.GrainSize = log.GrainSize }
	if log.Texture != "" { existingLog.Texture = log.Texture }
	if log.Structure != "" { existingLog.Structure = log.Structure }
	if log.Hardness != "" { existingLog.Hardness = log.Hardness }
	if log.Weathering != "" { existingLog.Weathering = log.Weathering }
	if log.Fracturing != "" { existingLog.Fracturing = log.Fracturing }
	if log.Description != "" { existingLog.Description = log.Description }
	if log.RQDPercentage.Internal.Value != "" { existingLog.RQDPercentage.Internal.Value = log.RQDPercentage.Internal.Value }
	if log.RecoveryPercentage.Internal.Value != "" { existingLog.RecoveryPercentage.Internal.Value = log.RecoveryPercentage.Internal.Value }
	if log.LoggedBy != "" { existingLog.LoggedBy = log.LoggedBy }
	if log.LoggedDate != nil && !log.LoggedDate.IsZero() { existingLog.LoggedDate = log.LoggedDate }

	existingLog.UpdatedAt = time.Now()

	appErr = s.lithologyRepo.UpdateLog(ctx, existingLog)
	if appErr != nil {
		return appErr
	}
	
	// --- LOG USER ACTIVITY ---
	newValueJSON, err := json.Marshal(existingLog)
	if err != nil {
		s.logger.Warn(ctx, "Failed to marshal new lithology log value", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "logID", Value: log.ID.String()})
	}
	newValueStr := string(newValueJSON)
	logIDStr := existingLog.ID.String()
	details := fmt.Sprintf("Lithology Log ID %s updated.", existingLog.ID.String())
	ipAddr := logCtx.IPAddress

	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, models.ActionTypeUpdateLithology, models.ResourceTypeLithology, &logIDStr, &ipAddr, &details, &oldValueStr, &newValueStr)
	
	if logErr != nil {
		s.logger.Error(ctx, "Failed to log UpdateLog activity", logErr, logger.Field{Key: "logID", Value: existingLog.ID.String()})
	}

	return nil
}

// DeleteLog handles the business logic for deleting a lithology log by its ID.
// REFACTORED: Fetches the log before deletion to capture its state.
func (s *LithologyServiceImpl) DeleteLog(ctx context.Context, id utils.BinaryUUID, userRole models.UserRole, logCtx models.ActivityLogContext) *exception.AppError {
	if userRole != models.RoleGeologist {
		return exception.NewPermissionError("User not authorized to delete lithology logs")
	}
	logToDelete, appErr := s.lithologyRepo.GetLogByID(ctx, id)
	if appErr != nil {
		return appErr
	}
	oldValueJSON, _ := json.Marshal(logToDelete)
	oldValueStr := string(oldValueJSON)

	appErr = s.lithologyRepo.DeleteLog(ctx, id)
	if appErr != nil {
		return appErr
	}
	
	// --- LOG USER ACTIVITY ---
	logIDStr := logToDelete.ID.String()
	details := fmt.Sprintf("Lithology Log ID %s deleted.", logToDelete.ID.String())
	ipAddr := logCtx.IPAddress

	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, models.ActionTypeDeleteLithology, models.ResourceTypeLithology, &logIDStr, &ipAddr, &details, &oldValueStr, nil)
	if logErr != nil {
		s.logger.Error(ctx, "Failed to log DeleteLog activity", logErr, logger.Field{Key: "logID", Value: logToDelete.ID.String()})
	}
	
	return nil
}

// ListLogsByStation retrieves a list of lithology logs associated with a specific station.
func (s *LithologyServiceImpl) ListLogsByStation(
	ctx context.Context,
	stationID utils.BinaryUUID,
) ([]*models.LithologyLog, *exception.AppError) {
	logs, appErr := s.lithologyRepo.ListLogsByStation(ctx, stationID)
	if appErr != nil {
		if appErr.Code == exception.ErrNotFound {
			return []*models.LithologyLog{}, nil
		}
		return nil, appErr
	}
	return logs, nil
}

// ListLogsByDepthRange retrieves lithology logs for a given station ID within a specified depth range.
func (s *LithologyServiceImpl) ListLogsByDepthRange(
	ctx context.Context,
	stationID utils.BinaryUUID,
	minDepth, maxDepth float64,
) ([]*models.LithologyLog, *exception.AppError) {
	logs, appErr := s.lithologyRepo.ListLogsByDepthRange(ctx, stationID, minDepth, maxDepth)
	if appErr != nil {
		if appErr.Code == exception.ErrNotFound {
			return []*models.LithologyLog{}, nil
		}
		return nil, appErr
	}
	return logs, nil
}