package contract

import (
	"boreholedata-ms/internal/api/dto"
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"context"
)

// StationRepository defines the contract for station data operations.
// StationRepository defines the contract for station data operations.
type StationRepository interface {
	Create(ctx context.Context, station *models.Station) *exception.AppError
	GetByID(ctx context.Context, id utils.BinaryUUID) (*models.Station, *exception.AppError)
	GetWithProject(ctx context.Context, id utils.BinaryUUID) (*models.Station, *exception.AppError)
	Update(ctx context.Context, station *models.Station) *exception.AppError
	Delete(ctx context.Context, id utils.BinaryUUID) *exception.AppError
	ListByProject(ctx context.Context, projectID utils.BinaryUUID) ([]*models.Station, *exception.AppError)
}

// StationService defines the contract for station business logic.
type StationService interface {
	CreateStation(ctx context.Context, station *models.Station, createdBy utils.BinaryUUID) (*models.Station, *exception.AppError)
	GetStation(ctx context.Context, id utils.BinaryUUID) (*models.Station, *exception.AppError)
	// --- FIX: Updated UpdateStation signature to include userRole ---
	UpdateStation(ctx context.Context, stationID utils.BinaryUUID, req *dto.UpdateStationRequest) (*models.Station, *exception.AppError)
	// --- END FIX ---
	DeleteStation(ctx context.Context, id utils.BinaryUUID) *exception.AppError
	ListStationsByProject(ctx context.Context, projectID utils.BinaryUUID) ([]*models.Station, *exception.AppError)
}
