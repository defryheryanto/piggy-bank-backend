package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/defryheryanto/piggy-bank-backend/internal/app"
	"github.com/defryheryanto/piggy-bank-backend/internal/auth"
	"github.com/defryheryanto/piggy-bank-backend/internal/httpserver/response"
)

func HandleRegister(a *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload *registerPayload

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			if err == io.EOF {
				response.WithJSON(w, http.StatusBadRequest, response.NewErrorResponse("Please fill form", "form is empty"))
				return
			}
			response.WithJSON(w, http.StatusInternalServerError, response.NewDefaultErrorResponse(err.Error()))
			return
		}
		err = payload.Validate()
		if err != nil {
			response.WithJSON(w, http.StatusInternalServerError, response.NewDefaultErrorResponse(err.Error()))
			return
		}

		user := &auth.User{
			Username: payload.Username,
			Password: payload.Password,
		}
		err = a.AuthService.Register(user)
		if err != nil {
			response.WithJSON(w, http.StatusInternalServerError, response.NewDefaultErrorResponse(err.Error()))
			return
		}

		response.WithJSON(w, http.StatusInternalServerError, response.NewSuccessResponse())
	}
}
