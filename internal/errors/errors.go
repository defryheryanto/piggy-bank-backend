package errors

import "net/http"

func NewBadRequestError(message, detail string) HandledError {
	return NewHandledError(http.StatusBadRequest, message, detail)
}
