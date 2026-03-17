package utils

import (
	"fmt"
)

const (
	ErrCodeNotFound     = "NOT_FOUND"
	ErrCodeDuplicate    = "DUPLICATE"
	ErrCodeValidation   = "VALIDATION"
	ErrCodeDatabase     = "DATABASE"
	ErrCodeConflict     = "CONFLICT"
	ErrCodeBadRequest   = "BAD_REQUEST"
	ErrCodeUnauthorized = "UNAUTHORIZED"
	ErrCodeForbidden    = "FORBIDDEN"
)

type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func WrapError(code, message string, err error) error {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func IsAppError(err error) (*AppError, bool) {
	if appErr, ok := err.(*AppError); ok {
		return appErr, true
	}
	return nil, false
}
