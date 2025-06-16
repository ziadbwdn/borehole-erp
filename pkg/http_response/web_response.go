package http_response

import (
	"boreholedata-ms/internal/exception" // Adjust import path if needed
	"net/http"

	"github.com/gin-gonic/gin"
)

// WebResponse is the standardized structure for all API responses.
type WebResponse struct {
	Code   int         `json:"code"`            // HTTP status code, e.g., 200
	Status string      `json:"status"`          // HTTP status text, e.g., "OK"
	Data   interface{} `json:"data,omitempty"`  // Payload for successful responses
	Error  interface{} `json:"error,omitempty"` // Structured error details for failed responses
}

// ErrorDetail is the standardized structure for the nested error object.
// This is created from an internal exception.AppError.
type ErrorDetail struct {
	Code    string      `json:"code"`              // Application-specific error code, e.g., "not_found"
	Message string      `json:"message"`           // Human-readable error message
	Details interface{} `json:"details,omitempty"` // Optional validation details or other info
}

// Success sends a standardized success response.
// It replaces the old RespondWithSuccess.
func RespondWithSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, WebResponse{
		Code:   code,
		Status: http.StatusText(code),
		Data:   data,
	})
}

// Error sends a standardized error response from an AppError.
// It replaces the old HandleAppError.
func HandleAppError(c *gin.Context, err *exception.AppError) {
	status := err.HTTPStatus()

	errorDetail := ErrorDetail{
		Code:    string(err.Code),
		Message: err.Message,
		Details: err.Details,
	}

	c.AbortWithStatusJSON(status, WebResponse{
		Code:   status,
		Status: http.StatusText(status),
		Error:  errorDetail,
	})
}
