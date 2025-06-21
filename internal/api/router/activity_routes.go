package router

import (
	"boreholedata-ms/internal/api/handler"
	"boreholedata-ms/internal/logger"
	"boreholedata-ms/internal/middleware"
	"boreholedata-ms/internal/models"

	"github.com/gin-gonic/gin"
)

func setupUserActivityRoutes(apiGroup *gin.RouterGroup, userActivityHandler *handler.UserActivityHandler, log logger.Logger) {
	activityGroup := apiGroup.Group("/activities")
	{
		activityGroup.POST("", userActivityHandler.LogActivity)
		activityGroup.GET("", middleware.AuthorizeRole(log, models.RoleAdmin), userActivityHandler.ListActivities)
		activityGroup.GET("/summary/:userID", middleware.AuthorizeRole(log, models.RoleAdmin), userActivityHandler.GetActivitySummary)
		activityGroup.GET("/alerts", middleware.AuthorizeRole(log, models.RoleAdmin), userActivityHandler.GetSecurityAlerts)
		activityGroup.DELETE("", middleware.AuthorizeRole(log, models.RoleAdmin), userActivityHandler.CleanOldActivities)
	}
}