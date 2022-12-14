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

func HandleGetUserConfig(a *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := auth.FromContext(r.Context())
		if session == nil {
			response.WithError(w, errors.ErrInvalidSession)
			return
		}

		cfg, err := a.UserConfigService.GetByUserId(session.UserID)
		if err != nil {
			response.WithError(w, err)
			return
		}

		response.WithJSON(w, http.StatusOK, cfg)
	}
}

func HandleUpdateUserConfig(a *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := &auth.UpdateUserConfigPayload{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			if err == io.EOF {
				response.WithError(w, errors.ErrEmptyPayload)
				return
			}
			response.WithError(w, errors.ErrUnprocessablePayload)
			return
		}

		session := auth.FromContext(r.Context())
		if session == nil {
			response.WithError(w, errors.ErrInvalidSession)
			return
		}

		payload.UserId = session.UserID

		err = a.UserConfigService.Update(payload)
		if err != nil {
			response.WithError(w, err)
			return
		}

		response.WithJSON(w, http.StatusOK, response.NewSuccessResponse())
	}
}
