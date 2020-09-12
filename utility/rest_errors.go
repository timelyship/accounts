package utility

import (
	"net/http"
	"time"
)

type RestError struct {
	Message   string `json:"Message`
	Status    int    `json:"Code`
	Error     string `json:"Error`
	Timestamp string `json:"Timestamp`
}

func NewBadRequestError(message string) *RestError {
	return &RestError{
		Message:   message,
		Status:    http.StatusBadRequest,
		Error:     "bad_request",
		Timestamp: InApiDateFormat(time.Now().UTC()),
	}
}

func NewInternalServerError(message string) *RestError {
	return &RestError{
		Message:   message,
		Status:    http.StatusInternalServerError,
		Error:     "interval_server_error",
		Timestamp: InApiDateFormat(time.Now().UTC()),
	}
}
