package services

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils" // Now importing utils for the helper functions
	"context"

	// "errors" // Removed: No longer directly used in this file
	"fmt"
	// "strconv" // No longer needed directly here
	"time"
	// pbdecimal "google.golang.org/genproto/googleapis/type/decimal" // No longer needed directly here
)

// StationServiceImpl implements the contract.StationService interface.
type StationServiceImpl struct {
	stationRepo contract.StationRepository
	projectRepo contract.ProjectRepository // Needed to verify ProjectID exists
}

// NewStationService creates and returns a new instance of StationServiceImpl.
func NewStationService(
	stationRepo contract.StationRepository,
	projectRepo contract.ProjectRepository, // Pass ProjectRepository dependency
) contract.StationService {
	return &StationServiceImpl{
		stationRepo: stationRepo,
		projectRepo: projectRepo,
	}
}

// CreateStation handles the business logic for creating a new station.
func (s *StationServiceImpl) CreateStation(
	ctx context.Context,
	station *models.Station,
	createdBy utils.BinaryUUID,
) (*models.Station, *exception.AppError) {
	// Verify ProjectID exists
	_, appErr := s.projectRepo.GetByID(ctx, station.ProjectID)
	if appErr != nil {
		// If project not found or other error, propagate it.
		// Note: If GetByID returns a *exception.AppError for NotFound, this will correctly propagate it.
		// If it returns a generic error, you might need errors.As here, but current contract says *AppError.
		return nil, exception.NewValidationError(fmt.Sprintf("Project with ID '%s' not found or inaccessible", station.ProjectID.String()))
	}

	station.ID = utils.NewBinaryUUID() // Generate a new UUID for the station
	// Set timestamps
	if station.CreatedAt.IsZero() {
		station.CreatedAt = time.Now()
	}
	station.UpdatedAt = time.Now()

	appErr = s.stationRepo.Create(ctx, station)
	if appErr != nil {
		return nil, appErr // Propagate error from repository
	}

	return station, nil
}

// GetStation retrieves a station by its ID.
func (s *StationServiceImpl) GetStation(
	ctx context.Context,
	id utils.BinaryUUID,
) (*models.Station, *exception.AppError) {
	station, appErr := s.stationRepo.GetByID(ctx, id)
	if appErr != nil {
		return nil, appErr // Propagate error from repository
	}
	return station, nil
}

// UpdateStation handles the business logic for updating an existing station.
func (s *StationServiceImpl) UpdateStation(
	ctx context.Context,
	station *models.Station, // This 'station' model should contain the ID and fields to update
) *exception.AppError {
	// First, retrieve the existing station to ensure it exists and to get current values.
	existingStation, appErr := s.stationRepo.GetByID(ctx, station.ID)
	if appErr != nil {
		return appErr // Propagate NotFoundError or DatabaseError from repo
	}

	// Apply updates from the provided 'station' model to the 'existingStation'.
	// Only update fields if they are explicitly provided (non-zero/non-empty for basic types, or checked for pointers in handler).
	// For decimal.Decimal, we need to check if the Value string is non-empty.
	if station.StationCode != "" {
		existingStation.StationCode = station.StationCode
	}
	if station.StationName != "" {
		existingStation.StationName = station.StationName
	}
	if station.StationType != "" {
		existingStation.StationType = station.StationType
	}
	// Corrected: Update the GormDecimal struct directly
	if station.Latitude.Internal.Value != "" { // Check if the string value is non-empty
		existingStation.Latitude.Internal.Value = station.Latitude.Internal.Value
	}
	if station.Longitude.Internal.Value != "" { // Check if the string value is non-empty
		existingStation.Longitude.Internal.Value = station.Longitude.Internal.Value
	}
	if station.Elevation.Internal.Value != "" { // Check if the string value is non-empty
		existingStation.Elevation.Internal.Value = station.Elevation.Internal.Value
	}
	if station.TotalDepth.Internal.Value != "" { // Check if the string value is non-empty
		existingStation.TotalDepth.Internal.Value = station.TotalDepth.Internal.Value
	}
	if !station.DrillingDate.IsZero() {
		existingStation.DrillingDate = station.DrillingDate
	}
	if station.GeologistName != "" {
		existingStation.GeologistName = station.GeologistName
	}
	if station.Notes != "" {
		existingStation.Notes = station.Notes
	}

	existingStation.UpdatedAt = time.Now() // Update the timestamp

	appErr = s.stationRepo.Update(ctx, existingStation)
	if appErr != nil {
		return appErr // Propagate error from repository
	}

	return nil
}

// DeleteStation handles the business logic for deleting a station by its ID.
func (s *StationServiceImpl) DeleteStation(
	ctx context.Context,
	id utils.BinaryUUID,
) *exception.AppError {
	appErr := s.stationRepo.Delete(ctx, id)
	if appErr != nil {
		return appErr // Propagate error from repository
	}
	return nil
}

// ListStationsByProject retrieves a list of stations associated with a specific project.
func (s *StationServiceImpl) ListStationsByProject(
	ctx context.Context,
	projectID utils.BinaryUUID,
) ([]*models.Station, *exception.AppError) {
	stations, appErr := s.stationRepo.ListByProject(ctx, projectID)
	if appErr != nil {
		// If no stations are found, the repo should return ErrNotFound.
		if appErr.Code == exception.ErrNotFound {
			return []*models.Station{}, nil // Return empty slice if no stations found
		}
		return nil, appErr // Propagate other errors
	}
	return stations, nil
}
