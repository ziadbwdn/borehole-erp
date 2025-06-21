package contract

import (
	"boreholedata-ms/internal/exception" // Import exception for AppError
	"boreholedata-ms/internal/utils"
	"boreholedata-ms/internal/models"
	"context"
)

// ReportService defines the contract for report generation business logic.
type ReportService interface {
	GenerateStationReport(ctx context.Context, stationID utils.BinaryUUID, logCtx models.ActivityLogContext) ([]byte, *exception.AppError)
	GenerateProjectSummary(ctx context.Context, projectID utils.BinaryUUID, logCtx models.ActivityLogContext) ([]byte, *exception.AppError)
}