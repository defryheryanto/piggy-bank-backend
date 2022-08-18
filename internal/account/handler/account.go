package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/defryheryanto/piggy-bank-backend/internal/account"
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

func HandleCreateAccount(a *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := auth.FromContext(r.Context())
		if session == nil {
			response.WithError(w, errors.InvalidSession)
			return
		}

		var payload *CreateAccountPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			if err == io.EOF {
				response.WithError(w, errors.EmptyPayload)
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

		acc := &account.Account{
			AccountTypeID: payload.AccountTypeId,
			AccountName:   payload.Name,
			UserID:        session.UserID,
		}

		err = a.AccountService.CreateAccount(acc)
		if err != nil {
			response.WithError(w, err)
			return
		}

		response.WithJSON(w, http.StatusOK, response.NewSuccessResponse())
	}
}
