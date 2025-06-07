package middleware

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/interfaces/contract" // Import the contract for AuthService
	"boreholedata-ms/pkg/http_response"            // For standardized error responses
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware struct holds the authentication service dependency.
type AuthMiddleware struct {
	authService contract.AuthService
}

// NewAuthMiddleware creates and returns a new instance of AuthMiddleware.
// It takes an AuthService implementation as a dependency.
func NewAuthMiddleware(authService contract.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// Handle returns a Gin middleware handler function.
// This function extracts the JWT token, validates it using the AuthService,
// and sets the userID and userRole in the Gin context for subsequent handlers.
func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			appErr := exception.NewAuthError("Authorization header required")
			http_response.HandleAppError(c, appErr)
			c.Abort() // Stop processing the request
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			appErr := exception.NewAuthError("Invalid authorization format. Expected 'Bearer <token>'")
			http_response.HandleAppError(c, appErr)
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Use the AuthService to verify the token
		userID, userRole, appErr := m.authService.VerifyToken(c.Request.Context(), tokenString)
		if appErr != nil {
			http_response.HandleAppError(c, appErr)
			c.Abort()
			return
		}

		// Set user context for downstream handlers
		c.Set("userID", userID)
		c.Set("userRole", userRole)
		c.Next() // Proceed to the next handler in the chain
	}
}
