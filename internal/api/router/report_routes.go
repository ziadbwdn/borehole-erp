package router

import (
	"boreholedata-ms/internal/api/handler"
	"boreholedata-ms/internal/logger"

	"github.com/gin-gonic/gin"
)

func setupReportRoutes(apiGroup *gin.RouterGroup, reportHandler *handler.ReportHandler, log logger.Logger) {
	reportGroup := apiGroup.Group("/reports")
	{
		reportGroup.GET("/stations/:station_id", reportHandler.GenerateStationReport)
		reportGroup.GET("/projects/:id/summary", reportHandler.GenerateProjectSummary)
	}
}