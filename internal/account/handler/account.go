package handler

import (
	"net/http"

	"github.com/defryheryanto/piggy-bank-backend/internal/app"
	"github.com/defryheryanto/piggy-bank-backend/internal/auth"
	"github.com/defryheryanto/piggy-bank-backend/internal/errors"
	"github.com/defryheryanto/piggy-bank-backend/internal/httpserver/response"
)

func HandleGetTypes(a *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		types := a.AccountService.GetTypes()
		response.WithJSON(w, http.StatusOK, types)
	}
}

func HandleGetAccounts(a *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := auth.FromContext(r.Context())
		if session == nil {
			response.WithError(w, errors.InvalidSession)
			return
		}

		accounts := a.AccountService.GetAccountsByUser(session.UserID)
		response.WithJSON(w, http.StatusOK, accounts)
	}
}
