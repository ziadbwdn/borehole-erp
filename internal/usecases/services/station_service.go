package services

import (
	"boreholedata-ms/internal/api/dto"
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"context"
	"fmt"
	"time"
)

// The service no longer needs a roleValidator.
type StationServiceImpl struct {
	stationRepo contract.StationRepository
	projectRepo contract.ProjectRepository
}

func NewStationService(
	stationRepo contract.StationRepository,
	projectRepo contract.ProjectRepository,
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
		if appErr.Code == exception.ErrNotFound {
			return nil, exception.NewValidationError(fmt.Sprintf("Project with ID '%s' not found or inaccessible", station.ProjectID.String()))
		}
		return nil, appErr // Propagate other errors from project repo
	}

	station.ID = utils.NewBinaryUUID() // Generate a new UUID for the station
	// Set timestamps
	if station.CreatedAt.IsZero() {
		station.CreatedAt = time.Now()
	}
	station.UpdatedAt = time.Now()

	// Ensure drilling status is valid if provided, otherwise default.
	// If the handler passes an empty string, it will be caught by the default GORM tag
	// (gorm:"default:'planned'"). However, if a non-empty, invalid string is passed,
	// validate it here.
	if station.DrillingStatus != "" {
		validStatus := false
		for _, s := range []models.DrillingStatus{
			models.DrillingStatusPlanned,
			models.DrillingStatusInProgress,
			models.DrillingStatusCompleted,
			models.DrillingStatusAbandoned,
			models.DrillingStatusSuspended,
			models.DrillingStatusOnHold,
		} {
			if station.DrillingStatus == s {
				validStatus = true
				break
			}
		}
		if !validStatus {
			return nil, exception.NewValidationError(fmt.Sprintf("Invalid drilling status for creation: %s", station.DrillingStatus))
		}
	}

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

// UpdateStation now contains all the logic and is much more robust.
func (s *StationServiceImpl) UpdateStation(
	ctx context.Context,
	stationID utils.BinaryUUID,
	req *dto.UpdateStationRequest, // Receives the DTO
) (*models.Station, *exception.AppError) {

	// Step 1: Load the full, existing station from the database.
	existingStation, appErr := s.stationRepo.GetByID(ctx, stationID)
	if appErr != nil {
		return nil, appErr
	}

	// Step 2: Apply changes directly from the DTO's non-nil fields.
	// This is now the single source of truth for update logic.
	if req.StationCode != nil {
		existingStation.StationCode = *req.StationCode
	}
	if req.StationName != nil {
		existingStation.StationName = *req.StationName
	}
	if req.StationType != nil {
		existingStation.StationType = *req.StationType
	}
	if req.DrillingStatus != nil {
		// You can add validation logic here if needed
		existingStation.DrillingStatus = models.DrillingStatus(*req.DrillingStatus)
	}
	if req.DrillingDate != nil {
		existingStation.DrillingDate = req.DrillingDate
	}
	if req.GeologistName != nil {
		existingStation.GeologistName = *req.GeologistName
	}
	if req.Notes != nil {
		existingStation.Notes = *req.Notes
	}

	// Handle GormDecimal fields
	if req.Latitude != nil {
		if gd, err := utils.StringToGormDecimal(*req.Latitude); err == nil {
			existingStation.Latitude = *gd
		}
	}
	if req.Longitude != nil {
		if gd, err := utils.StringToGormDecimal(*req.Longitude); err == nil {
			existingStation.Longitude = *gd
		}
	}
	if req.Elevation != nil {
		if gd, err := utils.StringToGormDecimal(*req.Elevation); err == nil {
			existingStation.Elevation = *gd
		}
	}
	if req.TotalDepth != nil {
		if gd, err := utils.StringToGormDecimal(*req.TotalDepth); err == nil {
			existingStation.TotalDepth = *gd
		}
	}
	// ... etc. for all fields in your DTO

	// Step 3: Save the modified record to the database.
	if appErr := s.stationRepo.Update(ctx, existingStation); appErr != nil {
		return nil, appErr
	}

	// Step 4: Return the fully updated model.
	return existingStation, nil
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
		if appErr.Code == exception.ErrNotFound {
			return []*models.Station{}, nil // Return empty slice if no stations found
		}
		return nil, appErr // Propagate other errors
	}
	return stations, nil
}
