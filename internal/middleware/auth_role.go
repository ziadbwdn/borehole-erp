package middleware

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/logger"
	"boreholedata-ms/internal/models" // Import models to use UserRole constants
	"boreholedata-ms/pkg/gin_helpers"
	"boreholedata-ms/pkg/http_response"

	"github.com/gin-gonic/gin"
)

// AuthorizeRole returns a Gin middleware that checks if the authenticated user's role
// is among the allowed roles for accessing a specific endpoint or resource.
func AuthorizeRole(log logger.Logger, allowedRoles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userRoleStr, appErr := gin_helpers.GetUserRoleFromContext(c)
		if appErr != nil {
			// This indicates an issue *retrieving* the role, possibly AuthMiddleware didn't run or failed.
			// This might still be an authentication (401) issue rather than permission (403).
			log.Error(ctx, "Authorization failed: Could not get user role from context", appErr)
			http_response.HandleAppError(c, exception.NewAuthError("Unauthorized: User role not found or invalid."))
			c.Abort()
			return
		}

		currentUserRole := models.UserRole(userRoleStr) // Cast the string from context to our UserRole type

		isAuthorized := false
		for _, role := range allowedRoles {
			if currentUserRole == role {
				isAuthorized = true
				break
			}
		}

		if !isAuthorized {
			log.Warn(ctx, "Authorization failed: User role is not permitted for this operation",
				logger.Field{Key: "userRole", Value: currentUserRole},
				logger.Field{Key: "allowedRoles", Value: allowedRoles})
			// Correctly using NewPermissionError which maps to ErrPermission -> HTTP 403 Forbidden
			http_response.HandleAppError(c, exception.NewPermissionError("Forbidden: You do not have the necessary permissions to perform this action."))
			c.Abort()
			return
		}

		c.Next() // User is authorized, proceed to the next handler in the chain
	}
}
