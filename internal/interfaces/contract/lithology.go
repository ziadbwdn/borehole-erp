package contract

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"context"
)

type LithologyRepository interface {
	CreateLog(ctx context.Context, log *models.LithologyLog) *exception.AppError
	GetLogByID(ctx context.Context, id utils.BinaryUUID) (*models.LithologyLog, *exception.AppError)
	UpdateLog(ctx context.Context, log *models.LithologyLog) *exception.AppError
	DeleteLog(ctx context.Context, id utils.BinaryUUID) *exception.AppError
	ListLogsByStation(ctx context.Context, stationID utils.BinaryUUID) ([]*models.LithologyLog, *exception.AppError)
	ListLogsByDepthRange(ctx context.Context, stationID utils.BinaryUUID, minDepth, maxDepth float64) ([]*models.LithologyLog, *exception.AppError)
}

// LithologyService defines the interface for lithology log related business operations.
type LithologyService interface {
	CreateLog(ctx context.Context, log *models.LithologyLog, createdBy utils.BinaryUUID, userRole models.UserRole) (*models.LithologyLog, *exception.AppError)
	GetLogByID(ctx context.Context, id utils.BinaryUUID) (*models.LithologyLog, *exception.AppError)
	UpdateLog(ctx context.Context, log *models.LithologyLog, userRole models.UserRole) *exception.AppError
	DeleteLog(ctx context.Context, id utils.BinaryUUID, userRole models.UserRole) *exception.AppError
	ListLogsByStation(ctx context.Context, stationID utils.BinaryUUID) ([]*models.LithologyLog, *exception.AppError)
	ListLogsByDepthRange(ctx context.Context, stationID utils.BinaryUUID, minDepth, maxDepth float64) ([]*models.LithologyLog, *exception.AppError)
}

