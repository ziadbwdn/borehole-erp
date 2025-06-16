package utils

import (
	"boreholedata-ms/internal/exception"
	"database/sql/driver"
	"fmt"
	"strconv"

	pbdecimal "google.golang.org/genproto/googleapis/type/decimal"
)

// GormDecimal is a wrapper around pbdecimal.Decimal that implements
// sql.Scanner and driver.Valuer for GORM compatibility.
// It holds pbdecimal.Decimal as a named field to avoid method/field name conflicts.
type GormDecimal struct {
	Internal pbdecimal.Decimal // Use a named field for the embedded struct
}

// Value implements the driver.Valuer interface for GormDecimal.
// It converts the GormDecimal to a string for database storage.
func (gd GormDecimal) Value() (driver.Value, error) {
	// Access the Value field of the internal pbdecimal.Decimal
	if gd.Internal.Value == "" {
		return nil, nil // Return nil for empty string to signify NULL in DB
	}
	// Return the string representation of the decimal value.
	return gd.Internal.Value, nil
}

// Scan implements the sql.Scanner interface for GormDecimal.
// It converts a database value (e.g., []byte, string, float64) back into a GormDecimal.
func (gd *GormDecimal) Scan(value interface{}) error {
	if value == nil {
		gd.Internal.Value = "0" // Set to "0" for null values to represent a zero decimal
		return nil
	}

	switch v := value.(type) {
	case []byte:
		gd.Internal.Value = string(v)
	case string:
		gd.Internal.Value = v
	case float64:
		// Convert float to string with sufficient precision
		gd.Internal.Value = strconv.FormatFloat(v, 'f', -1, 64) // -1 for dynamic precision
	case int64:
		gd.Internal.Value = strconv.FormatInt(v, 10)
	default:
		return fmt.Errorf("unsupported Scan type for GormDecimal: %T", value)
	}
	return nil
}

// StringToGormDecimal converts a string to a *utils.GormDecimal.
// This is the primary function to use when converting DTO strings to model decimals.
func StringToGormDecimal(s string) (*GormDecimal, *exception.AppError) {
	if s == "" {
		return &GormDecimal{Internal: pbdecimal.Decimal{Value: "0"}}, nil // Treat empty string as zero
	}
	// Basic validation: check if it's a valid number string.
	if _, err := strconv.ParseFloat(s, 64); err != nil {
		return nil, exception.NewValidationError(fmt.Sprintf("invalid decimal string format: '%s'", s), err.Error())
	}
	return &GormDecimal{Internal: pbdecimal.Decimal{Value: s}}, nil
}

// GormDecimalToString converts a *utils.GormDecimal to a string.
// This is the primary function to use when converting model decimals to DTO strings.
func GormDecimalToString(gd *GormDecimal) string {
	if gd == nil {
		return ""
	}
	return gd.Internal.Value
}
