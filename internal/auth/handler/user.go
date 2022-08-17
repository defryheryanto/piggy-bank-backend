package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/defryheryanto/piggy-bank-backend/internal/app"
	"github.com/defryheryanto/piggy-bank-backend/internal/auth"
	"github.com/defryheryanto/piggy-bank-backend/internal/errors"
	"github.com/defryheryanto/piggy-bank-backend/internal/httpserver/response"
)

func HandleRegister(a *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload *registerPayload

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			if err == io.EOF {
				response.WithError(w, errors.NewHandledError(http.StatusBadRequest, "Please fill form", "payload body is empty"))
				return
			}
			response.WithError(w, err)
			return
		}
		err = payload.Validate()
		if err != nil {
			response.WithError(w, err)
			return
		}

		user := &auth.User{
			Username: payload.Username,
			Password: payload.Password,
		}
		err = a.AuthService.Register(user)
		if err != nil {
			response.WithError(w, err)
			return
		}

		response.WithJSON(w, http.StatusInternalServerError, response.NewSuccessResponse())
	}
}
