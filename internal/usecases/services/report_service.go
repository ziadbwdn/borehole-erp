package services

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/internal/logger"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"boreholedata-ms/pkg/pdf"
	"context"
	"fmt"
)

type ReportServiceImpl struct {
	stationRepo     contract.StationRepository
	lithoRepo       contract.LithologyRepository
	labRepo         contract.LaboratoryRepository
	pdfGenerator    *pdf.PDFGenerator
	activityService contract.UserActivityService
	logger          logger.Logger
}

func NewReportService(stationRepo contract.StationRepository, lithoRepo contract.LithologyRepository, labRepo contract.LaboratoryRepository, activityService contract.UserActivityService, logger logger.Logger) contract.ReportService {
	if stationRepo == nil { panic("stationRepo must not be nil") }
	if lithoRepo == nil { panic("lithoRepo must not be nil") }
	if labRepo == nil { panic("labRepo must not be nil") }
	if activityService == nil { panic("activityService must not be nil") }
	if logger == nil { panic("logger must not be nil") }
	return &ReportServiceImpl{
		stationRepo:     stationRepo,
		lithoRepo:       lithoRepo,
		labRepo:         labRepo,
		pdfGenerator:    pdf.NewPDFGenerator(),
		activityService: activityService,
		logger:          logger,
	}
}

func (s *ReportServiceImpl) GenerateStationReport(ctx context.Context, stationID utils.BinaryUUID, logCtx models.ActivityLogContext) ([]byte, *exception.AppError) {
	station, appErr := s.stationRepo.GetByID(ctx, stationID)
	if appErr != nil { return nil, appErr }
	
	logs, appErr := s.lithoRepo.ListLogsByStation(ctx, stationID)
	if appErr != nil && appErr.Code != exception.ErrNotFound {
		return nil, exception.NewDatabaseError("lithology logs retrieval failed", appErr)
	}
	
	samples, appErr := s.labRepo.ListSamplesByStation(ctx, stationID)
	if appErr != nil && appErr.Code != exception.ErrNotFound {
		return nil, exception.NewDatabaseError("laboratory samples retrieval failed", appErr)
	}
	
	var ucsResults []*models.UCSResult
	for _, sample := range samples {
		results, appErrUCS := s.labRepo.GetUCSResultsBySample(ctx, sample.ID)
		if appErrUCS != nil && appErrUCS.Code != exception.ErrNotFound {
			s.logger.Warn(ctx, "Failed to get UCS results for sample during report generation", logger.Field{Key: "error", Value: appErrUCS.Error()})
			continue
		}
		if results != nil { ucsResults = append(ucsResults, results...) }
	}
	
	pdfBuffer, genErr := s.pdfGenerator.GenerateBoreholeLogPDF(station, logs, samples, ucsResults)
	if genErr != nil {
		return nil, exception.NewInternalError("PDF generation failed", genErr)
	}

	stationIDStr := station.ID.String()
	details := fmt.Sprintf("Generated station report for '%s'.", station.StationCode)
	ipAddr := logCtx.IPAddress
	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, models.ActionTypeGenerateReport, models.ResourceTypeReport, &stationIDStr, &ipAddr, &details, nil, nil)
	if logErr != nil {
		s.logger.Error(ctx, "Failed to log GenerateStationReport activity", logErr, logger.Field{Key: "stationID", Value: station.ID.String()})
	}

	return pdfBuffer.Bytes(), nil
}

func (s *ReportServiceImpl) GenerateProjectSummary(ctx context.Context, projectID utils.BinaryUUID, logCtx models.ActivityLogContext) ([]byte, *exception.AppError) {
	projectIDStr := projectID.String()
	details := fmt.Sprintf("Generated project summary report for Project ID %s.", projectID.String())
	ipAddr := logCtx.IPAddress
	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, models.ActionTypeGenerateReport, models.ResourceTypeReport, &projectIDStr, &ipAddr, &details, nil, nil)
	if logErr != nil {
		s.logger.Error(ctx, "Failed to log GenerateProjectSummary activity", logErr, logger.Field{Key: "projectID", Value: projectID.String()})
	}

	return nil, exception.NewInternalError("GenerateProjectSummary not yet implemented", nil)
}