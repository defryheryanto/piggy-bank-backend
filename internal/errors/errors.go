package errors

import "net/http"

func NewBadRequestError(message, detail string) HandledError {
	return NewHandledError(http.StatusBadRequest, message, detail)
}

func NewUnauthorizedError(message, detail string) HandledError {
	return NewHandledError(http.StatusUnauthorized, message, detail)
}

func NewNotFoundError(message, detail string) HandledError {
	return NewHandledError(http.StatusNotFound, message, detail)
}

func NewUnprocessableEntity(message, detail string) HandledError {
	return NewHandledError(http.StatusUnprocessableEntity, message, detail)
}

var InvalidSession = NewUnauthorizedError("Unauthorized", "invalid session")
var EmptyPayload = NewBadRequestError("Please fill data", "request body is empty")
var ErrUnprocessablePayload = NewUnprocessableEntity("Data submitted invalid", "payload is unprocessable")
