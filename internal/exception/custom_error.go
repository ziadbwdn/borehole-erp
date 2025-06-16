package exception

import (
	"fmt"
	"net/http"
	"strings"
)

type ErrorCode string

const (
	ErrDatabase      ErrorCode = "database_error"
	ErrNotFound      ErrorCode = "not_found"
	ErrValidation    ErrorCode = "validation_failed"
	ErrAuth          ErrorCode = "authentication_error"
	ErrPermission    ErrorCode = "permission_denied"
	ErrPDFGeneration ErrorCode = "pdf_generation_failed"
	ErrInternal      ErrorCode = "internal_error"
)

type AppError struct {
	Code      ErrorCode `json:"code"`
	Message   string    `json:"message"`
	Details   []string  `json:"details,omitempty"`
	Operation string    `json:"operation,omitempty"`
	Err       error     `json:"-"`
}

func (e *AppError) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] %s", e.Code, e.Message))

	if len(e.Details) > 0 {
		sb.WriteString(" | Details: ")
		sb.WriteString(strings.Join(e.Details, "; "))
	}

	if e.Err != nil {
		sb.WriteString(fmt.Sprintf(" | Cause: %v", e.Err))
	}

	return sb.String()
}

func (e *AppError) HTTPStatus() int {
	switch e.Code {
	case ErrDatabase:
		return http.StatusServiceUnavailable
	case ErrNotFound:
		return http.StatusNotFound
	case ErrValidation:
		return http.StatusBadRequest
	case ErrAuth:
		return http.StatusUnauthorized
	case ErrPermission:
		return http.StatusForbidden
	case ErrPDFGeneration:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// Constructors

// NewDatabaseError creates an AppError for database-related issues.
func NewDatabaseError(op string, err error) *AppError {
	return &AppError{
		Code:      ErrDatabase,
		Message:   "Database operation failed",
		Operation: op,
		Err:       err,
	}
}

// NewNotFoundError creates an AppError for resources not found.
func NewNotFoundError(entity string, id interface{}) *AppError {
	return &AppError{
		Code:    ErrNotFound,
		Message: fmt.Sprintf("%s not found", entity),
		Details: []string{fmt.Sprintf("identifier: %v", id)},
	}
}

// NewValidationError creates an AppError for invalid input.
func NewValidationError(msg string, details ...string) *AppError {
	return &AppError{
		Code:    ErrValidation,
		Message: msg,
		Details: details,
	}
}

// NewAuthError creates an AppError for authentication failures (e.g., invalid token, missing credentials).
func NewAuthError(msg string, details ...string) *AppError {
	return &AppError{
		Code:    ErrAuth,
		Message: msg,
		Details: details,
	}
}

// NewPermissionError creates an AppError for authorization/permission failures (HTTP 403 Forbidden).
func NewPermissionError(msg string) *AppError {
	return &AppError{
		Code:    ErrPermission,
		Message: msg,
	}
}

// NewInternalError creates an AppError for unexpected internal server issues.
func NewInternalError(op string, err error) *AppError {
	return &AppError{
		Code:      ErrInternal,
		Message:   "Internal server error",
		Operation: op,
		Err:       err,
	}
}

// NewPDFGenerationError creates an AppError for PDF generation failures.
func NewPDFGenerationError(op string, err error) *AppError {
	return &AppError{
		Code:      ErrPDFGeneration,
		Message:   "PDF generation failed",
		Operation: op,
		Err:       err,
	}
}
