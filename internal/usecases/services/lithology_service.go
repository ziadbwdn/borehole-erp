package services

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils" // For BinaryUUID and decimal conversions
	"context"
	"fmt"
	"time"
)

// LithologyServiceImpl implements the contract.LithologyService interface.
type LithologyServiceImpl struct {
	lithologyRepo contract.LithologyRepository
	stationRepo   contract.StationRepository // Dependency to verify StationID exists
}

// NewLithologyService creates and returns a new instance of LithologyServiceImpl.
func NewLithologyService(
	lithologyRepo contract.LithologyRepository,
	stationRepo contract.StationRepository,
) contract.LithologyService {
	return &LithologyServiceImpl{
		lithologyRepo: lithologyRepo,
		stationRepo:   stationRepo,
	}
}

// CreateLog handles the business logic for creating a new lithology log.
func (s *LithologyServiceImpl) CreateLog(
	ctx context.Context,
	log *models.LithologyLog,
	createdBy utils.BinaryUUID, // This parameter is currently unused for the log model itself
	userRole models.UserRole, // Added userRole for authorization
) (*models.LithologyLog, *exception.AppError) {
	// --- Authorization Check for Create ---
	// Only 'admin' or 'geologist' roles can create lithology logs.
	if userRole != models.RoleAdmin && userRole != models.RoleGeologist {
		// Using NewPermissionError for authorization issues
		return nil, exception.NewPermissionError("User not authorized to create lithology logs")
	}
	// --- End Authorization Check ---

	// Verify StationID exists
	_, appErr := s.stationRepo.GetByID(ctx, log.StationID)
	if appErr != nil {
		return nil, exception.NewValidationError(fmt.Sprintf("Station with ID '%s' not found or inaccessible", log.StationID.String()))
	}

	// Set new UUID and timestamps
	log.ID = utils.NewBinaryUUID()
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}
	log.UpdatedAt = time.Now()
	// If LithologyLog had a CreatedByID field, you would set it here:
	// log.CreatedByID = createdBy

	appErr = s.lithologyRepo.CreateLog(ctx, log)
	if appErr != nil {
		return nil, appErr // Propagate error from repository
	}

	return log, nil
}

// GetLogByID retrieves a lithology log by its ID.
// No specific role-based access check here, assuming public read access or that
// project-level access is handled implicitly by the client's scope.
func (s *LithologyServiceImpl) GetLogByID(
	ctx context.Context,
	id utils.BinaryUUID,
) (*models.LithologyLog, *exception.AppError) {
	log, appErr := s.lithologyRepo.GetLogByID(ctx, id)
	if appErr != nil {
		return nil, appErr // Propagate error from repository
	}
	return log, nil
}

// UpdateLog handles the business logic for updating an existing lithology log.
func (s *LithologyServiceImpl) UpdateLog(
	ctx context.Context,
	log *models.LithologyLog,
	userRole models.UserRole, // Added userRole for authorization
) *exception.AppError {
	// --- Authorization Check for Update ---
	// Only 'admin' or 'geologist' roles can update lithology logs.
	// For more granular control (e.g., only update their own logs),
	// you would fetch the existing log and compare 'loggedBy' field with the current user's ID.
	if userRole != models.RoleAdmin && userRole != models.RoleGeologist {
		// Using NewPermissionError for authorization issues
		return exception.NewPermissionError("User not authorized to update lithology logs")
	}
	// --- End Authorization Check ---

	// First, retrieve the existing log to ensure it exists and to get current values.
	existingLog, appErr := s.lithologyRepo.GetLogByID(ctx, log.ID)
	if appErr != nil {
		return appErr // Propagate NotFoundError or DatabaseError from repo
	}

	// Apply updates from the provided 'log' model to the 'existingLog'.
	// Only update fields if they are explicitly provided (non-empty for strings, non-zero for time, non-empty for decimal.Decimal.Value).
	if log.StationID != (utils.BinaryUUID{}) { // Check if StationID is provided (non-zero UUID)
		// Verify new StationID exists if it's being changed
		_, appErr := s.stationRepo.GetByID(ctx, log.StationID)
		if appErr != nil {
			return exception.NewValidationError(fmt.Sprintf("New Station with ID '%s' not found or inaccessible", log.StationID.String()))
		}
		existingLog.StationID = log.StationID
	}
	if log.DepthFrom.Internal.Value != "" {
		existingLog.DepthFrom.Internal.Value = log.DepthFrom.Internal.Value
	}
	if log.DepthTo.Internal.Value != "" {
		existingLog.DepthTo.Internal.Value = log.DepthTo.Internal.Value
	}
	if log.LithologyType != "" {
		existingLog.LithologyType = log.LithologyType
	}
	if log.RockColor != "" {
		existingLog.RockColor = log.RockColor
	}
	if log.GrainSize != "" {
		existingLog.GrainSize = log.GrainSize
	}
	if log.Texture != "" {
		existingLog.Texture = log.Texture
	}
	if log.Structure != "" {
		existingLog.Structure = log.Structure
	}
	if log.Hardness != "" {
		existingLog.Hardness = log.Hardness
	}
	if log.Weathering != "" {
		existingLog.Weathering = log.Weathering
	}
	if log.Fracturing != "" {
		existingLog.Fracturing = log.Fracturing
	}
	if log.Description != "" {
		existingLog.Description = log.Description
	}
	// IMPORTANT: Check if pbdecimal.Decimal.Value is not the zero-value string "0"
	// if it's meant to distinguish between "not provided" and "explicitly zero".
	// For now, assuming any non-empty string means it's provided.
	if log.RQDPercentage.Internal.Value != "" {
		existingLog.RQDPercentage.Internal.Value = log.RQDPercentage.Internal.Value
	}
	if log.RecoveryPercentage.Internal.Value != "" {
		existingLog.RecoveryPercentage.Internal.Value = log.RecoveryPercentage.Internal.Value
	}
	if log.LoggedBy != "" {
		existingLog.LoggedBy = log.LoggedBy
	}
	if log.LoggedDate != nil && !log.LoggedDate.IsZero() { // Ensure pointer is not nil before dereferencing and checking IsZero()
		existingLog.LoggedDate = log.LoggedDate
	}

	existingLog.UpdatedAt = time.Now() // Update the timestamp

	appErr = s.lithologyRepo.UpdateLog(ctx, existingLog)
	if appErr != nil {
		return appErr // Propagate error from repository
	}

	return nil
}

// DeleteLog handles the business logic for deleting a lithology log by its ID.
func (s *LithologyServiceImpl) DeleteLog(
	ctx context.Context,
	id utils.BinaryUUID,
	userRole models.UserRole, // Added userRole for authorization
) *exception.AppError {
	// --- Authorization Check for Delete ---
	// Only 'admin' role can delete lithology logs.
	// You might allow 'editor'/'geologist' to delete their own, but that requires
	// fetching the log first and comparing the 'createdBy' or 'loggedBy' field.
	if userRole != models.RoleAdmin {
		// Using NewPermissionError for authorization issues
		return exception.NewPermissionError("User not authorized to delete lithology logs")
	}
	// --- End Authorization Check ---

	appErr := s.lithologyRepo.DeleteLog(ctx, id)
	if appErr != nil {
		return appErr // Propagate error from repository
	}
	return nil
}

// ListLogsByStation retrieves a list of lithology logs associated with a specific station.
// No specific role-based access check here, assuming public read access or that
// project-level access is handled implicitly by the client's scope.
func (s *LithologyServiceImpl) ListLogsByStation(
	ctx context.Context,
	stationID utils.BinaryUUID,
) ([]*models.LithologyLog, *exception.AppError) {
	logs, appErr := s.lithologyRepo.ListLogsByStation(ctx, stationID)
	if appErr != nil {
		if appErr.Code == exception.ErrNotFound {
			return []*models.LithologyLog{}, nil // Return empty slice if no logs found
		}
		return nil, appErr // Propagate other errors
	}
	return logs, nil
}

// ListLogsByDepthRange retrieves lithology logs for a given station ID within a specified depth range.
// No specific role-based access check here, similar to ListLogsByStation.
func (s *LithologyServiceImpl) ListLogsByDepthRange(
	ctx context.Context,
	stationID utils.BinaryUUID,
	minDepth, maxDepth float64, // These are float64 as per the interface
) ([]*models.LithologyLog, *exception.AppError) {
	logs, appErr := s.lithologyRepo.ListLogsByDepthRange(ctx, stationID, minDepth, maxDepth)
	if appErr != nil {
		if appErr.Code == exception.ErrNotFound {
			return []*models.LithologyLog{}, nil // Return empty slice if no logs found
		}
		return nil, appErr // Propagate other errors
	}
	return logs, nil
}
