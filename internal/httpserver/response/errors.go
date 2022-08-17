package response

import (
	"encoding/json"
	"net/http"

	"github.com/defryheryanto/piggy-bank-backend/internal/errors"
)

func WithError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	definedError := defineError(err)
	w.WriteHeader(definedError.HttpStatus)
	json.NewEncoder(w).Encode(definedError)
}

func defineError(err error) errors.HandledError {
	var definedError errors.HandledError
	switch err := err.(type) {
	case errors.HandledError:
		definedError = err
	default:
		definedError = errors.NewDefaultHandledError(http.StatusInternalServerError, err.Error())
	}

	return definedError
}
