package gin_helpers

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/utils" // Correctly imports the 'utils' package from internal/utils
	"boreholedata-ms/pkg/http_response"
	"fmt"

	"github.com/gin-gonic/gin"
)

// ParseIDFromContext is a utility function to extract and parse a UUID from Gin context parameters.
func ParseIDFromContext(c *gin.Context, paramName, resourceType string) (utils.BinaryUUID, *exception.AppError) {
	idStr := c.Param(paramName)
	if idStr == "" {
		appErr := exception.NewValidationError(fmt.Sprintf("%s ID is required", resourceType), fmt.Sprintf("Missing path parameter '%s'", paramName))
		http_response.HandleAppError(c, appErr)
		return utils.BinaryUUID{}, appErr
	}
	id, err := utils.ParseBinaryUUID(idStr)
	if err != nil {
		appErr := exception.NewValidationError(fmt.Sprintf("Invalid %s ID format", resourceType), err.Error())
		http_response.HandleAppError(c, appErr)
		return utils.BinaryUUID{}, appErr
	}
	return id, nil
}
