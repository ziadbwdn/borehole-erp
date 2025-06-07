package repository // This MUST be "repository"

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"context"
	"fmt"

	"gorm.io/gorm"
)

// stationRepository is a concrete implementation of the contract.StationRepository interface.
type stationRepository struct {
	db *gorm.DB
}

// NewStationRepository creates a new instance of StationRepository.
func NewStationRepository(db *gorm.DB) contract.StationRepository {
	return &stationRepository{db: db}
}

// Create inserts a new station record into the database.
func (r *stationRepository) Create(ctx context.Context, station *models.Station) *exception.AppError {
	if err := r.db.WithContext(ctx).Create(station).Error; err != nil {
		return exception.NewDatabaseError("Failed to create station", err)
	}
	return nil
}

// GetByID retrieves a single station record by its ID.
func (r *stationRepository) GetByID(ctx context.Context, id utils.BinaryUUID) (*models.Station, *exception.AppError) {
	var station models.Station
	if err := r.db.WithContext(ctx).First(&station, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundError("Station", id.String())
		}
		return nil, exception.NewDatabaseError(fmt.Sprintf("Failed to retrieve station with ID %s", id.String()), err)
	}
	return &station, nil
}

// GetWithProject retrieves a station record by its ID, preloading its associated project.
func (r *stationRepository) GetWithProject(ctx context.Context, id utils.BinaryUUID) (*models.Station, *exception.AppError) {
	var station models.Station
	if err := r.db.WithContext(ctx).Preload("Project").First(&station, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundError("Station", id.String())
		}
		return nil, exception.NewDatabaseError(fmt.Sprintf("Failed to retrieve station with ID %s", id.String()), err)
	}
	return &station, nil
}

// Update updates an existing station record in the database.
func (r *stationRepository) Update(ctx context.Context, station *models.Station) *exception.AppError {
	if err := r.db.WithContext(ctx).Save(station).Error; err != nil {
		return exception.NewDatabaseError(fmt.Sprintf("Failed to update station with ID %s", station.ID.String()), err)
	}
	return nil
}

// Delete deletes a station record by its ID.
func (r *stationRepository) Delete(ctx context.Context, id utils.BinaryUUID) *exception.AppError {
	if err := r.db.WithContext(ctx).Delete(&models.Station{}, "id = ?", id).Error; err != nil {
		return exception.NewDatabaseError(fmt.Sprintf("Failed to delete station with ID %s", id.String()), err)
	}
	return nil
}

// ListByProject retrieves all stations associated with a given project ID.
func (r *stationRepository) ListByProject(ctx context.Context, projectID utils.BinaryUUID) ([]*models.Station, *exception.AppError) {
	var stations []*models.Station
	if err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Find(&stations).Error; err != nil {
		return nil, exception.NewDatabaseError(fmt.Sprintf("Failed to retrieve stations for project ID %s", projectID.String()), err)
	}
	if len(stations) == 0 {
		return nil, exception.NewNotFoundError("Stations", fmt.Sprintf("for project ID %s", projectID.String()))
	}
	return stations, nil
}
