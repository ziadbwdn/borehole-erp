package contract

import (
	"boreholedata-ms/internal/exception" // Import exception for AppError
	"boreholedata-ms/internal/utils"
	"context"
)

// ReportService defines the contract for report generation business logic.
type ReportService interface {
	// Corrected return types to match the service implementation and handler's needs.
	// It now returns a byte slice for the PDF content and our custom AppError.
	GenerateStationReport(ctx context.Context, stationID utils.BinaryUUID) ([]byte, *exception.AppError)
	GenerateProjectSummary(ctx context.Context, projectID utils.BinaryUUID) ([]byte, *exception.AppError)
}
