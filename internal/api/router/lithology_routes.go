package router

import (
	"boreholedata-ms/internal/api/handler"
	"boreholedata-ms/internal/logger"
	"boreholedata-ms/internal/middleware"
	"boreholedata-ms/internal/models"

	"github.com/gin-gonic/gin"
)

func setupLithologyRoutes(apiGroup *gin.RouterGroup, lithologyHandler *handler.LithologyHandler, log logger.Logger) {
	lithologyGroup := apiGroup.Group("/lithology-logs")
	{
		lithologyGroup.POST("", middleware.AuthorizeRole(log, models.RoleGeologist), lithologyHandler.CreateLithologyLog)
		lithologyGroup.PUT("/:id", middleware.AuthorizeRole(log, models.RoleGeologist), lithologyHandler.UpdateLithologyLog)
		lithologyGroup.DELETE("/:id", middleware.AuthorizeRole(log, models.RoleGeologist), lithologyHandler.DeleteLithologyLog)
		lithologyGroup.GET("/:id", lithologyHandler.GetLithologyLogByID)
	}
}