package services

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"boreholedata-ms/pkg/pdf"
	"context"
	"fmt"
)

type ReportServiceImpl struct {
	stationRepo contract.StationRepository
	lithoRepo   contract.LithologyRepository
	labRepo     contract.LaboratoryRepository
	pdfGenerator *pdf.PDFGenerator
}

func NewReportService(
	stationRepo contract.StationRepository,
	lithoRepo contract.LithologyRepository,
	labRepo contract.LaboratoryRepository,
) *ReportServiceImpl {
	return &ReportServiceImpl{
		stationRepo: stationRepo,
		lithoRepo:   lithoRepo,
		labRepo:     labRepo,
		pdfGenerator: pdf.NewPDFGenerator(),
	}
}

func (s *ReportServiceImpl) GenerateStationReport(
	ctx context.Context,
	stationID utils.BinaryUUID,
) ([]byte, *exception.AppError) {
	// Get station (without project relation for now)
	station, appErr := s.stationRepo.GetByID(ctx, stationID) // Assuming GetByID returns *exception.AppError
	if appErr != nil {
		// Propagate the error directly from the repository
		return nil, appErr
	}

	// Get lithology logs
	logs, appErr := s.lithoRepo.ListLogsByStation(ctx, stationID) // Assuming ListLogsByStation returns *exception.AppError
	if appErr != nil {
		if appErr.Code == exception.ErrNotFound {
			logs = []*models.LithologyLog{} // Treat as empty if not found
		} else {
			return nil, exception.NewDatabaseError("lithology logs retrieval failed", appErr)
		}
	}

	// Get lab samples
	samples, appErr := s.labRepo.ListSamplesByStation(ctx, stationID) // Assuming ListSamplesByStation returns *exception.AppError
	if appErr != nil {
		if appErr.Code == exception.ErrNotFound {
			samples = []*models.LabSample{} // Treat as empty if not found
		} else {
			return nil, exception.NewDatabaseError("laboratory samples retrieval failed", appErr)
		}
	}

	// Get UCS results
	var ucsResults []*models.UCSResult
	for _, sample := range samples {
		results, appErrUCS := s.labRepo.GetUCSResultsBySample(ctx, sample.ID) // Assuming GetUCSResultsBySample returns *exception.AppError
		if appErrUCS != nil {
			if appErrUCS.Code == exception.ErrNotFound {
				continue // If no UCS results found for this sample, continue to the next.
			}
			// It's a genuine AppError (not NotFound). Log it and continue.
			fmt.Printf("Warning: Failed to get UCS results for sample %s: %v\n", sample.ID.String(), appErrUCS)
			continue
		}
		ucsResults = append(ucsResults, results...)
	}

	// Generate PDF report
	generator := pdf.NewPDFGenerator()
	// Corrected: Assign to a generic 'error' variable first
	pdfBuffer, genErr := generator.GenerateBoreholeLogPDF(station, logs, samples, ucsResults)
	if genErr != nil {
		// Corrected: Pass the underlying error (genErr) as the second argument
		return nil, exception.NewInternalError("PDF generation failed", genErr)
	}

	return pdfBuffer.Bytes(), nil
}

// GenerateProjectSummary will generate a summary report for a given project.
// This is a placeholder for future implementation.
func (s *ReportServiceImpl) GenerateProjectSummary(ctx context.Context, projectID utils.BinaryUUID) ([]byte, *exception.AppError) {
	// Corrected: Use exception.NewInternalError with the correct signature
	return nil, exception.NewInternalError("GenerateProjectSummary not yet implemented", nil) // Pass nil for the error if no underlying error
}
