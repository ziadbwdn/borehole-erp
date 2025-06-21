package contract

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"context"
)

// LaboratoryRepository defines the contract for laboratory data operations.
type LaboratoryRepository interface {
	CreateSample(ctx context.Context, sample *models.LabSample) *exception.AppError
	GetSampleByID(ctx context.Context, id utils.BinaryUUID) (*models.LabSample, *exception.AppError)
	UpdateSample(ctx context.Context, sample *models.LabSample) *exception.AppError
	DeleteSample(ctx context.Context, id utils.BinaryUUID) *exception.AppError
	ListSamplesByStation(ctx context.Context, stationID utils.BinaryUUID) ([]*models.LabSample, *exception.AppError)

	CreateUCSResult(ctx context.Context, ucs *models.UCSResult) *exception.AppError
	GetUCSResultByID(ctx context.Context, id utils.BinaryUUID) (*models.UCSResult, *exception.AppError)
	UpdateUCSResult(ctx context.Context, ucs *models.UCSResult) *exception.AppError
	DeleteUCSResult(ctx context.Context, id utils.BinaryUUID) *exception.AppError
	// Corrected: Method name to match the service's call and return type for multiple results
	GetUCSResultsBySample(ctx context.Context, sampleID utils.BinaryUUID) ([]*models.UCSResult, *exception.AppError)
}

// LaboratoryService defines the contract for laboratory business logic.
type LaboratoryService interface {
	CreateSample(ctx context.Context, sample *models.LabSample, userRole models.UserRole, logCtx models.ActivityLogContext) (*models.LabSample, *exception.AppError)
	GetSample(ctx context.Context, id utils.BinaryUUID) (*models.LabSample, *exception.AppError)
	UpdateSample(ctx context.Context, sample *models.LabSample, userRole models.UserRole, logCtx models.ActivityLogContext) *exception.AppError
	DeleteSample(ctx context.Context, id utils.BinaryUUID, userRole models.UserRole, logCtx models.ActivityLogContext) *exception.AppError
	ListSamplesByStation(ctx context.Context, stationID utils.BinaryUUID) ([]*models.LabSample, *exception.AppError)

	CreateUCSResult(ctx context.Context, ucs *models.UCSResult, userRole models.UserRole, logCtx models.ActivityLogContext) (*models.UCSResult, *exception.AppError)
	GetUCSResult(ctx context.Context, id utils.BinaryUUID) (*models.UCSResult, *exception.AppError)
	UpdateUCSResult(ctx context.Context, ucs *models.UCSResult, userRole models.UserRole, logCtx models.ActivityLogContext) *exception.AppError
	DeleteUCSResult(ctx context.Context, id utils.BinaryUUID, userRole models.UserRole, logCtx models.ActivityLogContext) *exception.AppError
	ListUCSResultsBySample(ctx context.Context, sampleID utils.BinaryUUID) ([]*models.UCSResult, *exception.AppError)
}
