package router

import (
	"boreholedata-ms/internal/api/handler"
	"boreholedata-ms/internal/middleware"

	"github.com/gin-gonic/gin"
)

func setupAuthRoutes(authGroup *gin.RouterGroup, authHandler *handler.AuthHandler, authMiddleware *middleware.AuthMiddleware) {
	// Public routes
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/refresh", authHandler.Refresh)

	authGroup.POST("/logout", authHandler.Logout)

	// Protected profile routes
	profileGroup := authGroup.Group("/profile")
	profileGroup.Use(authMiddleware.Handle())
	{
		profileGroup.GET("", authHandler.GetProfile)
		profileGroup.PUT("", authHandler.UpdateProfile)
	}
}