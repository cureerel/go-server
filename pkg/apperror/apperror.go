// pkg/apperror/apperror.go
package apperror

import (
	"errors"
	"net/http"
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}
func (e *AppError) Unwrap() error { return e.Err }
func (e *AppError) HTTPStatus() int {
	if e.Code == 0 {
		return http.StatusInternalServerError
	}
	return e.Code
}
func (e *AppError) PublicError() string {
	if e.Message != "" {
		return e.Message
	}
	return "something went wrong"
}

// Free functions used by handlers
func HTTPStatus(err error) int {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae.HTTPStatus()
	}
	return http.StatusInternalServerError
}

func PublicMessage(err error) string {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae.PublicError()
	}
	return "something went wrong"
}

// Constructors
func NewBadRequest(msg string) *AppError   { return &AppError{Code: 400, Message: msg} }
func NewUnauthorized(msg string) *AppError { return &AppError{Code: 401, Message: msg} }
func NewForbidden(msg string) *AppError    { return &AppError{Code: 403, Message: msg} }
func NewNotFound(msg string) *AppError     { return &AppError{Code: 404, Message: msg} }
func NewConflict(msg string) *AppError     { return &AppError{Code: 409, Message: msg} }
func NewValidation(field, msg string) *AppError {
	return &AppError{Code: 422, Message: field + ": " + msg}
}
func NewInternal(err error, msg string) *AppError {
	return &AppError{Code: 500, Message: msg, Err: err}
}
func Wrap(err error, code int, msg string) *AppError {
	return &AppError{Code: code, Message: msg, Err: err}
}

// Type checks
func IsNotFound(err error) bool     { return codeIs(err, 404) }
func IsBadRequest(err error) bool   { return codeIs(err, 400) }
func IsConflict(err error) bool     { return codeIs(err, 409) }
func IsForbidden(err error) bool    { return codeIs(err, 403) }
func IsUnauthorized(err error) bool { return codeIs(err, 401) }

func codeIs(err error, code int) bool {
	var ae *AppError
	return errors.As(err, &ae) && ae.Code == code
}
