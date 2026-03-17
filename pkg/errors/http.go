package errors

import (
	"encoding/json"
	"net/http"
)

type HTTPResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewSuccessResponse(message string) *SuccessResponse {
	return &SuccessResponse{
		Success: true,
		Message: message,
	}
}

func NewHTTPResponse(err AppError) *HTTPResponse {
	return &HTTPResponse{
		Code:    err.Code,
		Message: err.Message,
		Data:    nil,
	}
}

func NewHTTPResponseWithData(data interface{}) *HTTPResponse {
	return &HTTPResponse{
		Code:    "",
		Message: "",
		Data:    data,
	}
}

func WriteHTTPError(w http.ResponseWriter, err AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(getHTTPStatusCode(err.Code))
	json.NewEncoder(w).Encode(NewHTTPResponse(err))
}

func WriteHTTPSuccess(w http.ResponseWriter, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewSuccessResponse(message))
}

func WriteHTTPSuccessWithData(w http.ResponseWriter, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewSuccessResponse(message))
}

func getHTTPStatusCode(code string) int {
	switch code {
	case CodeBadRequest:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeConflict:
		return http.StatusConflict
	case CodeValidationError:
		return http.StatusUnprocessableEntity
	case CodeDatabaseError:
		return http.StatusInternalServerError
	case CodeInternalError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
