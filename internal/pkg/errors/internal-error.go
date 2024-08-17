package errors

import "net/http"

type internalAPIError struct {
	message string
	status  int
}

func (e *internalAPIError) Error() string {
	return e.message
}

func (e *internalAPIError) StatusCode() int {
	return e.status
}

func NewInternalError(message string) APIError {
	return &internalAPIError{message: message, status: http.StatusInternalServerError}
}

func NewValidationError(message string) APIError {
	return &internalAPIError{message: message, status: http.StatusBadRequest}
}

func NewNotFoundError(message string) APIError {
	return &internalAPIError{message: message, status: http.StatusNotFound}
}

func NewAuthError(message string) APIError {
	return &internalAPIError{message: message, status: http.StatusUnauthorized}
}
