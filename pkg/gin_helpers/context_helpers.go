package gin_helpers

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext extracts the userID from Gin's context.
// It returns the BinaryUUID and an AppError if the userID is not found or is in an invalid format.
func GetUserIDFromContext(c *gin.Context) (utils.BinaryUUID, *exception.AppError) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		return utils.BinaryUUID{}, exception.NewAuthError("User ID not found in context")
	}
	userID, ok := userIDAny.(utils.BinaryUUID)
	if !ok {
		return utils.BinaryUUID{}, exception.NewInternalError("Invalid user ID format in context", nil)
	}
	return userID, nil
}

// GetUserRoleFromContext extracts the user role from Gin's context.
// It returns the role string and an AppError if the role is not found or is in an invalid format.
func GetUserRoleFromContext(c *gin.Context) (string, *exception.AppError) {
	userRoleAny, exists := c.Get("userRole") // Assuming "userRole" is the key for role
	if !exists {
		return "", exception.NewAuthError("User role not found in context")
	}
	userRole, ok := userRoleAny.(string)
	if !ok {
		return "", exception.NewInternalError("Invalid user role format in context", nil)
	}
	return userRole, nil
}
