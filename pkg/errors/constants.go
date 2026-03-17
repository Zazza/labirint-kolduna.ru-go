package errors

const (
	CodeInternalError   = "INTERNAL_ERROR"
	CodeBadRequest      = "BAD_REQUEST"
	CodeUnauthorized    = "UNAUTHORIZED"
	CodeForbidden       = "FORBIDDEN"
	CodeNotFound        = "NOT_FOUND"
	CodeConflict        = "CONFLICT"
	CodeValidationError = "VALIDATION_ERROR"
	CodeDatabaseError   = "DATABASE_ERROR"
)

var (
	ErrInternal     = NewAppError(CodeInternalError, "internal server error")
	ErrBadRequest   = NewAppError(CodeBadRequest, "bad request")
	ErrUnauthorized = NewAppError(CodeUnauthorized, "unauthorized")
	ErrForbidden    = NewAppError(CodeForbidden, "forbidden")
	ErrNotFound     = NewAppError(CodeNotFound, "resource not found")
	ErrConflict     = NewAppError(CodeConflict, "resource conflict")
	ErrValidation   = NewAppError(CodeValidationError, "validation failed")
	ErrDatabase     = NewAppError(CodeDatabaseError, "database error")
)

type AppError struct {
	Code    string
	Message string
	Err     error
}

func NewAppError(code string, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     nil,
	}
}

func NewAppErrorWithErr(code string, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
