package errors

import "net/http"

type AppError struct {
    Code int
    Err  error
}

func (e *AppError) Error() string {
    return e.Err.Error()
}

func NewBadRequest(err error) *AppError {
    return &AppError{Code: http.StatusBadRequest, Err: err}
}

func NewInternal(err error) *AppError {
    return &AppError{Code: http.StatusInternalServerError, Err: err}
}

// Optional: Helper to create a generic AppError easily
func New(code int, message string) *AppError {
    return &AppError{Code: code, Err: &customErr{message}}
}

type customErr struct {
    msg string
}

func (e *customErr) Error() string {
    return e.msg
}