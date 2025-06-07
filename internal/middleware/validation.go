package middleware

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ValidateRequest is a generic validation middleware
func ValidateRequest[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T

		// Bind JSON
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(exception.NewValidationError("invalid request format"))
			c.Abort()
			return
		}

		// Validate struct
		if err := utils.ValidateStruct(req); err != nil {
			validationErrors := utils.GetValidationErrors(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": validationErrors,
			})
			c.Abort()
			return
		}

		// Store validated request in context
		c.Set("validatedRequest", req)
		c.Next()
	}
}
