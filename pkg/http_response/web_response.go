package http_response

import (
	"boreholedata-ms/internal/api/dto"
	"boreholedata-ms/internal/exception"

	"github.com/gin-gonic/gin"
)

// HandleAppError is a utility function to standardize error responses from AppError.
// It takes a Gin context and an AppError, and sends a JSON response with the appropriate
// HTTP status code and error details.
func HandleAppError(c *gin.Context, appErr *exception.AppError) {
	c.JSON(appErr.HTTPStatus(), dto.ErrorResponse{
		Code:    string(appErr.Code),
		Message: appErr.Message,
		Details: appErr.Details,
	})
}

// RespondWithSuccess is a utility function to standardize success responses.
// It takes a Gin context, an HTTP status code, and the data to be returned.
func RespondWithSuccess(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}
