package utils

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
)

func init() {
	validate = validator.New()

	// Register custom validations
	_ = validate.RegisterValidation("password", validatePassword)
	_ = validate.RegisterValidation("username", validateUsername)
	_ = validate.RegisterValidation("geolocation", validateGeoLocation)
}

// ValidateStruct performs validation using centralized rules
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// Custom validation functions
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Minimum 12 characters
	if len(password) < 12 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	// Alphanumeric with underscores, 3-50 characters
	return regexp.MustCompile(`^[a-zA-Z0-9_]{3,50}$`).MatchString(username)
}

func validateGeoLocation(fl validator.FieldLevel) bool {
	coords := fl.Field().String()
	// Validate latitude/longitude format
	return regexp.MustCompile(`^-?\d{1,3}\.\d{1,6},\s*-?\d{1,3}\.\d{1,6}$`).MatchString(coords)
}

// ValidationError simplifies error extraction
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func GetValidationErrors(err error) []ValidationError {
	var ve []ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		ve = append(ve, ValidationError{
			Field:   err.Field(),
			Message: validationMessage(err),
		})
	}
	return ve
}

func validationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "password":
		return "Password must be at least 12 characters with uppercase, lowercase, number, and special character"
	case "username":
		return "Username must be 3-50 alphanumeric characters"
	case "geolocation":
		return "Invalid geographic coordinates format"
	case "min":
		return "Value too short"
	case "max":
		return "Value too long"
	default:
		return "Validation failed"
	}
}
