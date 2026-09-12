package errors

import (
	stdErrors "errors"
	"fmt"
	"net/http"
)

type Code string

const (
	CodeInvalidInput     Code = "INVALID_INPUT"
	CodeNotFound         Code = "NOT_FOUND"
	CodeInternal         Code = "INTERNAL_ERROR"
	CodeTimeout          Code = "TIMEOUT"
	CodeDependencyFailed Code = "DEPENDENCY_FAILED"
)

type AppError struct {
	Code      Code        `json:"code"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	Retryable bool        `json:"retryable"`
	Status    int         `json:"-"`
	Cause     error       `json:"-"`
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func InvalidInput(message string, details interface{}) *AppError {
	return &AppError{Code: CodeInvalidInput, Message: message, Details: details, Retryable: false, Status: http.StatusBadRequest}
}

func NotFound(message string, details interface{}) *AppError {
	return &AppError{Code: CodeNotFound, Message: message, Details: details, Retryable: false, Status: http.StatusNotFound}
}

func Internal(message string, cause error) *AppError {
	return &AppError{Code: CodeInternal, Message: message, Retryable: false, Status: http.StatusInternalServerError, Cause: cause}
}

func Timeout(message string, cause error) *AppError {
	return &AppError{Code: CodeTimeout, Message: message, Retryable: true, Status: http.StatusGatewayTimeout, Cause: cause}
}

func DependencyFailed(message string, cause error, retryable bool) *AppError {
	return &AppError{Code: CodeDependencyFailed, Message: message, Retryable: retryable, Status: http.StatusBadGateway, Cause: cause}
}

func AsAppError(err error) *AppError {
	if err == nil {
		return nil
	}
	var appErr *AppError
	if stdErrors.As(err, &appErr) {
		return appErr
	}
	return Internal("unexpected error", err)
}

func StatusCode(err error) int {
	return AsAppError(err).Status
}
