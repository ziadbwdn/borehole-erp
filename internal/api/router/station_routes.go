package router

import (
	"boreholedata-ms/internal/api/handler"
	"boreholedata-ms/internal/logger"
	"boreholedata-ms/internal/middleware"
	"boreholedata-ms/internal/models"

	"github.com/gin-gonic/gin"
)

func setupStationRoutes(apiGroup *gin.RouterGroup, stationHandler *handler.StationHandler, lithologyHandler *handler.LithologyHandler, laboratoryHandler *handler.LaboratoryHandler, log logger.Logger) {
	stationGroup := apiGroup.Group("/stations")
	{
		stationGroup.POST("", middleware.AuthorizeRole(log, models.RoleEngineer, models.RoleGeologist), stationHandler.CreateStation)
		stationGroup.PUT("/:id", middleware.AuthorizeRole(log, models.RoleEngineer, models.RoleGeologist), stationHandler.UpdateStation)
		stationGroup.DELETE("/:id", middleware.AuthorizeRole(log, models.RoleEngineer, models.RoleGeologist), stationHandler.DeleteStation)
		stationGroup.GET("/:id", stationHandler.GetStation)
		stationGroup.GET("/:id/lithology-logs", lithologyHandler.ListLithologyLogsByStation)
		stationGroup.GET("/:id/lithology-logs/by-depth", lithologyHandler.ListLithologyLogsByDepthRange)
		stationGroup.GET("/:id/lab-samples", laboratoryHandler.ListLabSamplesByStation)
	}
}