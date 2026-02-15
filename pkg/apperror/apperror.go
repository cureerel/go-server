package apperror

import (
    "errors"
    "net/http"
)

// AppError is our custom error type
type AppError struct {
    Code    int    // HTTP status code
    Message string // User-friendly message
    Err     error  // Internal error (for logs)
}

// Error implements the error interface
func (e *AppError) Error() string {
    if e.Err != nil {
        return e.Err.Error()
    }
    return e.Message
}

// Unwrap allows errors.Is and errors.As to work
func (e *AppError) Unwrap() error {
    return e.Err
}

// HTTPStatus returns the HTTP status code
func (e *AppError) HTTPStatus() int {
    if e.Code == 0 {
        return http.StatusInternalServerError
    }
    return e.Code
}

// PublicError returns safe message for users
func (e *AppError) PublicError() string {
    if e.Message != "" {
        return e.Message
    }
    return "Something went wrong"
}

// ==================== CONSTRUCTORS ====================

func NewBadRequest(message string) *AppError {
    return &AppError{
        Code:    http.StatusBadRequest,
        Message: message,
    }
}

func NewNotFound(message string) *AppError {
    return &AppError{
        Code:    http.StatusNotFound,
        Message: message,
    }
}

func NewUnauthorized(message string) *AppError {
    return &AppError{
        Code:    http.StatusUnauthorized,
        Message: message,
    }
}

func NewForbidden(message string) *AppError {
    return &AppError{
        Code:    http.StatusForbidden,
        Message: message,
    }
}

func NewInternal(err error, message string) *AppError {
    return &AppError{
        Code:    http.StatusInternalServerError,
        Message: message,
        Err:     err,
    }
}

func NewValidation(field, message string) *AppError {
    return &AppError{
        Code:    http.StatusUnprocessableEntity, // 422
        Message: field + ": " + message,
    }
}

// ==================== WRAPPER ====================

func Wrap(err error, code int, message string) *AppError {
    return &AppError{
        Code:    code,
        Message: message,
        Err:     err,
    }
}

// ==================== CHECK FUNCTIONS ====================

func IsNotFound(err error) bool {
    var appErr *AppError
    if errors.As(err, &appErr) {
        return appErr.Code == http.StatusNotFound
    }
    return false
}

func IsBadRequest(err error) bool {
    var appErr *AppError
    if errors.As(err, &appErr) {
        return appErr.Code == http.StatusBadRequest
    }
    return false
}