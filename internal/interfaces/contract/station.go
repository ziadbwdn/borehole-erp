package contract

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"context"
)

// StationRepository defines the contract for station data operations.
type StationRepository interface {
	Create(ctx context.Context, station *models.Station) *exception.AppError                 // Changed error type
	GetByID(ctx context.Context, id utils.BinaryUUID) (*models.Station, *exception.AppError) // Changed error type
	// GetWithProject is already *exception.AppError
	GetWithProject(ctx context.Context, id utils.BinaryUUID) (*models.Station, *exception.AppError)
	Update(ctx context.Context, station *models.Station) *exception.AppError                                // Changed error type
	Delete(ctx context.Context, id utils.BinaryUUID) *exception.AppError                                    // Changed error type
	ListByProject(ctx context.Context, projectID utils.BinaryUUID) ([]*models.Station, *exception.AppError) // Changed error type
}

// StationService defines the contract for station business logic.
type StationService interface {
	// Corrected CreateStation signature to match implementation
	CreateStation(ctx context.Context, station *models.Station, createdBy utils.BinaryUUID) (*models.Station, *exception.AppError)
	GetStation(ctx context.Context, id utils.BinaryUUID) (*models.Station, *exception.AppError)                     // Changed error type
	UpdateStation(ctx context.Context, station *models.Station) *exception.AppError                                 // Changed error type
	DeleteStation(ctx context.Context, id utils.BinaryUUID) *exception.AppError                                     // Changed error type
	ListStationsByProject(ctx context.Context, projectID utils.BinaryUUID) ([]*models.Station, *exception.AppError) // Changed error type
}
