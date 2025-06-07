package gin_helpers

import (
	"boreholedata-ms/internal/exception"
	"boreholedata-ms/internal/utils" // Correctly imports the 'utils' package from internal/utils
	"boreholedata-ms/pkg/http_response"
	"fmt"

	"github.com/gin-gonic/gin"
)

// ParseIDFromContext is a utility function to extract and parse a UUID from Gin context parameters.
// It takes the Gin context, the name of the URL parameter (e.g., "id", "station_id"),
// and a descriptive name for the entity (e.g., "lab sample", "station").
// It returns the parsed BinaryUUID and an *exception.AppError if parsing fails.
// If an error occurs, it also handles sending the appropriate HTTP response via http_response.HandleAppError.
func ParseIDFromContext(c *gin.Context, paramName, entityName string) (utils.BinaryUUID, *exception.AppError) {
	idStr := c.Param(paramName)
	// Calling the exported ParseBinaryUUID from the imported 'utils' package
	id, err := utils.ParseBinaryUUID(idStr) // This line should now resolve correctly
	if err != nil {
		appErr := exception.NewValidationError(fmt.Sprintf("Invalid %s ID format", entityName), err.Error())
		http_response.HandleAppError(c, appErr)
		return utils.BinaryUUID{}, appErr // Return zero UUID and the error
	}
	return id, nil
}
