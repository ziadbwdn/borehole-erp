package handler

import (
	"boreholedata-ms/internal/interfaces/contract" // Assuming contract.ReportService is here
	"boreholedata-ms/internal/models"
	"boreholedata-ms/pkg/gin_helpers"
	"boreholedata-ms/pkg/http_response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ReportHandler handles HTTP requests related to report generation.
type ReportHandler struct {
	reportService contract.ReportService
	authService   contract.AuthService
}

// NewReportHandler creates and returns a new instance of ReportHandler.
func NewReportHandler(reportService contract.ReportService, authService contract.AuthService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
		authService:   authService,
	}
}

// NewReportHandler creates and returns a new instance of ReportHandler.
func (h *ReportHandler) GenerateStationReport(c *gin.Context) {
	stationID, appErr := gin_helpers.ParseIDFromContext(c, "station_id", "station")
	if appErr != nil {
		// ParseIDFromContext already handles the response
		return
	}

	// 1. Get UserID and Username for logging
	userID, appErr := gin_helpers.GetUserIDFromContext(c); if appErr != nil { http_response.HandleAppError(c, appErr); return }
	username, appErr := h.authService.GetUserDetailsForLogging(c.Request.Context(), userID); if appErr != nil { http_response.HandleAppError(c, appErr); return }

	// 2. Prepare the ActivityLogContext
	logCtx := models.ActivityLogContext{
		UserID:    userID.String(),
		Username:  username,
		IPAddress: c.ClientIP(),
	}

	// 3. Call the service with the new logCtx parameter
	pdfBytes, appErr := h.reportService.GenerateStationReport(c.Request.Context(), stationID, logCtx)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=station_report.pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// GenerateProjectSummary handles the request to generate a summary PDF report for a specific project.
// @Router /api/reports/project/{project_id} [get]
func (h *ReportHandler) GenerateProjectSummary(c *gin.Context) {
	projectID, appErr := gin_helpers.ParseIDFromContext(c, "id", "project")
	if appErr != nil {
		return
	}

	// 1. Get UserID and Username for logging
	userID, appErr := gin_helpers.GetUserIDFromContext(c); if appErr != nil { http_response.HandleAppError(c, appErr); return }
	username, appErr := h.authService.GetUserDetailsForLogging(c.Request.Context(), userID); if appErr != nil { http_response.HandleAppError(c, appErr); return }

	// 2. Prepare the ActivityLogContext
	logCtx := models.ActivityLogContext{
		UserID:    userID.String(),
		Username:  username,
		IPAddress: c.ClientIP(),
	}

	// 3. Call the service with the new logCtx parameter
	pdfBytes, appErr := h.reportService.GenerateProjectSummary(c.Request.Context(), projectID, logCtx)
	if appErr != nil {
		http_response.HandleAppError(c, appErr)
		return
	}
	
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=project_summary.pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}