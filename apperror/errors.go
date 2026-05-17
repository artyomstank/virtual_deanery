// apperror/errors.go
package apperror

import (
	"fmt"
)

// Code represents application error code.
type Code string

const (
	// NOT_FOUND when resource not found
	NOT_FOUND Code = "NOT_FOUND"

	// CONFLICT when resource already exists
	CONFLICT Code = "CONFLICT"

	// BAD_REQUEST when input is invalid
	BAD_REQUEST Code = "BAD_REQUEST"

	// INVALID_CREDENTIALS when auth fails
	INVALID_CREDENTIALS Code = "INVALID_CREDENTIALS"

	// FORBIDDEN when user lacks permissions
	FORBIDDEN Code = "FORBIDDEN"

	// INTERNAL_SERVER for unexpected errors
	INTERNAL_SERVER Code = "INTERNAL_SERVER"

	// VALIDATION_ERROR for validation failures
	VALIDATION_ERROR Code = "VALIDATION_ERROR"

	// UNAUTHENTICATED when token is missing/invalid
	UNAUTHENTICATED Code = "UNAUTHENTICATED"
)

// AppError represents application error.
type AppError struct {
	Code    Code
	Message string
	Cause   error
}

// Error implements error interface.
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Is checks if error matches code.
func (e *AppError) Is(code Code) bool {
	return e.Code == code
}

// Sentinel errors
var (
	ErrNotFound           = &AppError{Code: NOT_FOUND, Message: "resource not found"}
	ErrConflict           = &AppError{Code: CONFLICT, Message: "resource already exists"}
	ErrBadRequest         = &AppError{Code: BAD_REQUEST, Message: "invalid request"}
	ErrInvalidCredentials = &AppError{Code: INVALID_CREDENTIALS, Message: "invalid credentials"}
	ErrForbidden          = &AppError{Code: FORBIDDEN, Message: "access denied"}
	ErrInternalServer     = &AppError{Code: INTERNAL_SERVER, Message: "internal server error"}
	ErrValidation         = &AppError{Code: VALIDATION_ERROR, Message: "validation failed"}
	ErrUnauthenticated    = &AppError{Code: UNAUTHENTICATED, Message: "authentication required"}
)

// New creates new app error.
func New(code Code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps error with app error.
func Wrap(code Code, message string, cause error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// IsAppError checks if error is AppError.
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetCode extracts error code from error.
func GetCode(err error) Code {
	if ae, ok := err.(*AppError); ok {
		return ae.Code
	}
	return INTERNAL_SERVER
}
