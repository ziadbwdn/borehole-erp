package gin_helpers

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext extracts the userID from Gin's context.
// It returns the BinaryUUID and an AppError if the userID is not found or is in an invalid format.

// Define context keys
const (
	UserIDContextKey   = "userID"
	UserRoleContextKey = "userRole" // New key for user role
)

func GetUserIDFromContext(c *gin.Context) (utils.BinaryUUID, *exception.AppError) {
	val, exists := c.Get(UserIDContextKey)
	if !exists {
		appErr := exception.NewAuthError("User ID not found in context. Authentication middleware missing or failed.")
		c.AbortWithStatus(http.StatusUnauthorized) // Or Internal Server Error if unexpected
		return utils.BinaryUUID{}, appErr
	}
	userID, ok := val.(utils.BinaryUUID)
	if !ok {
		appErr := exception.NewInternalError("Invalid user ID type in context.", nil)
		c.AbortWithStatus(http.StatusInternalServerError)
		return utils.BinaryUUID{}, appErr
	}
	return userID, nil
}

// GetUserRoleFromContext extracts the user role from Gin's context.
// It returns the role string and an AppError if the role is not found or is in an invalid format.
func GetUserRoleFromContext(c *gin.Context) (models.UserRole, *exception.AppError) {
	val, exists := c.Get(UserRoleContextKey)
	if !exists {
		appErr := exception.NewAuthError("User role not found in context. Authentication middleware missing or failed.")
		c.AbortWithStatus(http.StatusUnauthorized)
		return "", appErr
	}
	userRoleStr, ok := val.(string) // Role is stored as string by AuthMiddleware
	if !ok {
		appErr := exception.NewInternalError("Invalid user role type in context.", nil)
		c.AbortWithStatus(http.StatusInternalServerError)
		return "", appErr
	}
	userRole := models.UserRole(userRoleStr) // Cast string back to UserRole
	return userRole, nil
}
