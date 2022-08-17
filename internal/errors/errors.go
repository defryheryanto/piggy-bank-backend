package errors

import "net/http"

func NewBadRequestError(message, detail string) HandledError {
	return NewHandledError(http.StatusBadRequest, message, detail)
}

func NewUnauthorizedError(message, detail string) HandledError {
	return NewHandledError(http.StatusUnauthorized, message, detail)
}

var InvalidSession = NewUnauthorizedError("Unauthorized", "invalid session")
