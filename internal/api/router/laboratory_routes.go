package router

import (
	"boreholedata-ms/internal/api/handler"
	"boreholedata-ms/internal/logger"
	"boreholedata-ms/internal/middleware"
	"boreholedata-ms/internal/models"

	"github.com/gin-gonic/gin"
)

func setupLabSampleRoutes(apiGroup *gin.RouterGroup, laboratoryHandler *handler.LaboratoryHandler, log logger.Logger) {
	labSampleGroup := apiGroup.Group("/lab-samples")
	{
		labSampleGroup.POST("", middleware.AuthorizeRole(log, models.RoleLabTechnician), laboratoryHandler.CreateLabSample)
		labSampleGroup.PUT("/:id", middleware.AuthorizeRole(log, models.RoleLabTechnician), laboratoryHandler.UpdateLabSample)
		labSampleGroup.DELETE("/:id", middleware.AuthorizeRole(log, models.RoleLabTechnician), laboratoryHandler.DeleteLabSample)
		labSampleGroup.GET("/:id", laboratoryHandler.GetLabSampleByID)
		labSampleGroup.GET("/:id/ucs-results", laboratoryHandler.ListUCSResultsBySample)
	}
}

func setupUCSResultRoutes(apiGroup *gin.RouterGroup, laboratoryHandler *handler.LaboratoryHandler, log logger.Logger) {
	ucsResultGroup := apiGroup.Group("/ucs-results")
	{
		ucsResultGroup.POST("", middleware.AuthorizeRole(log, models.RoleLabTechnician), laboratoryHandler.CreateUCSResult)
		ucsResultGroup.PUT("/:id", middleware.AuthorizeRole(log, models.RoleLabTechnician), laboratoryHandler.UpdateUCSResult)
		ucsResultGroup.DELETE("/:id", middleware.AuthorizeRole(log, models.RoleLabTechnician), laboratoryHandler.DeleteUCSResult)
		ucsResultGroup.GET("/:id", laboratoryHandler.GetUCSResultByID)
	}
}