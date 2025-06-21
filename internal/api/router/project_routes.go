package router

import (
	"boreholedata-ms/internal/api/handler"
	"boreholedata-ms/internal/logger"
	"boreholedata-ms/internal/middleware"
	"boreholedata-ms/internal/models"

	"github.com/gin-gonic/gin"
)

func setupProjectRoutes(apiGroup *gin.RouterGroup, projectHandler *handler.ProjectHandler, stationHandler *handler.StationHandler, log logger.Logger) {
	projectGroup := apiGroup.Group("/projects")
	{
		projectGroup.POST("", middleware.AuthorizeRole(log, models.RoleEngineer), projectHandler.CreateProject)
		projectGroup.PUT("/:id", middleware.AuthorizeRole(log, models.RoleEngineer), projectHandler.UpdateProject)
		projectGroup.DELETE("/:id", middleware.AuthorizeRole(log, models.RoleEngineer), projectHandler.DeleteProject)
		projectGroup.GET("", projectHandler.ListProjects)
		projectGroup.GET("/:id", projectHandler.GetProject)
		projectGroup.GET("/:id/stations", stationHandler.ListStationsByProject)
		projectGroup.GET("/:id/stations/export/geo", stationHandler.ExportStationPointsGeo)
		projectGroup.GET("/:id/stations/export/utm", stationHandler.ExportStationPointsUTM)
	}
}