package utility

import (
	"errors"
	"net/http"
	"os"
	"time"
	"timelyship.com/accounts/application"
)

type RestError struct {
	Message   string `json:"message"`
	Status    int    `json:"code"`
	Error     error  `json:"error"`
	Timestamp string `json:"timestamp"`
}

func NewBadRequestError(message string, err *error) *RestError {
	return &RestError{
		Message:   message,
		Status:    http.StatusBadRequest,
		Error:     maskedError(err, "bad_request"),
		Timestamp: InApiDateFormat(time.Now().UTC()),
	}
}

func NewInternalServerError(message string, err *error) *RestError {
	return &RestError{
		Message:   message,
		Status:    http.StatusInternalServerError,
		Error:     maskedError(err, "interval_server_error"),
		Timestamp: InApiDateFormat(time.Now().UTC()),
	}
}

func NewUnAuthorizedError(message string, err *error) *RestError {
	return &RestError{
		Message:   message,
		Status:    http.StatusUnauthorized,
		Error:     maskedError(err, "unauthorized"),
		Timestamp: InApiDateFormat(time.Now().UTC()),
	}
}

func NewForbiddenError(message string, err *error) *RestError {
	return &RestError{
		Message:   message,
		Status:    http.StatusForbidden,
		Error:     maskedError(err, "bad_request"),
		Timestamp: InApiDateFormat(time.Now().UTC()),
	}
}

/*private methods*/
func maskedError(e *error, s string) error {
	// applies mask if log level is info
	if e == nil || os.Getenv("LOG_LEVEL") == application.STRING_CONST.LOG_LEVEL_INFO {
		return errors.New(s)
	}
	// otherwise does not apply mask
	return *e
}
