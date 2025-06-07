package handler

import (
	"boreholedata-ms/internal/interfaces/contract" // Assuming contract.ReportService is here
	"boreholedata-ms/pkg/gin_helpers"
	"boreholedata-ms/pkg/http_response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ReportHandler handles HTTP requests related to report generation.
type ReportHandler struct {
	reportService contract.ReportService // Dependency on the report service
}

// NewReportHandler creates and returns a new instance of ReportHandler.
func NewReportHandler(reportService contract.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

// GenerateStationReport handles the request to generate a PDF report for a specific station.
// @Router /api/reports/station/{station_id} [get]
func (h *ReportHandler) GenerateStationReport(c *gin.Context) {
	stationID, appErr := gin_helpers.ParseIDFromContext(c, "station_id", "station")
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	pdfBytes, appErr := h.reportService.GenerateStationReport(c.Request.Context(), stationID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Set headers for PDF download
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=station_report.pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// GenerateProjectSummary handles the request to generate a summary PDF report for a specific project.
// @Router /api/reports/project/{project_id} [get]
func (h *ReportHandler) GenerateProjectSummary(c *gin.Context) {
	// Corrected: Parse 'id' from context as per the new route definition
	projectID, appErr := gin_helpers.ParseIDFromContext(c, "id", "project")
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// Call the service method (which currently returns "Not Implemented")
	pdfBytes, appErr := h.reportService.GenerateProjectSummary(c.Request.Context(), projectID)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	// This part will only be reached if GenerateProjectSummary actually returns a PDF.
	// For now, it's a placeholder.
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=project_summary.pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
