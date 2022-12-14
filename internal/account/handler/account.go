package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/defryheryanto/piggy-bank-backend/internal/account"
	"github.com/defryheryanto/piggy-bank-backend/internal/app"
	"github.com/defryheryanto/piggy-bank-backend/internal/auth"
	"github.com/defryheryanto/piggy-bank-backend/internal/errors"
	"github.com/defryheryanto/piggy-bank-backend/internal/httpserver/response"
	"github.com/gorilla/mux"
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
			response.WithError(w, errors.ErrInvalidSession)
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
			response.WithError(w, errors.ErrInvalidSession)
			return
		}

		var payload *CreateAccountPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			if err == io.EOF {
				response.WithError(w, errors.ErrEmptyPayload)
				return
			}
			response.WithError(w, errors.ErrUnprocessablePayload)
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

func HandleUpdateAccount(a *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountIdStr := mux.Vars(r)["account_id"]
		accountId, err := strconv.Atoi(accountIdStr)
		if err != nil {
			response.WithError(w, errors.NewBadRequestError("Account ID should be numeric", "account_id should be numeric"))
			return
		}
		var p *account.UpdateAccountPayload

		err = json.NewDecoder(r.Body).Decode(&p)
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
		p.UserID = session.UserID
		p.AccountID = accountId

		err = a.AccountService.UpdateAccount(p)
		if err != nil {
			response.WithError(w, err)
			return
		}

		response.WithJSON(w, http.StatusOK, response.NewSuccessResponse())
	}
}

func HandleDeleteAccount(a *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountIdStr := mux.Vars(r)["account_id"]
		accountId, err := strconv.Atoi(accountIdStr)
		if err != nil {
			response.WithError(w, errors.NewBadRequestError("Account ID should be numeric", "account_id should be numeric"))
			return
		}

		session := auth.FromContext(r.Context())
		if session == nil {
			response.WithError(w, errors.ErrInvalidSession)
			return
		}

		err = a.AccountService.DeleteById(accountId, session.UserID)
		if err != nil {
			response.WithError(w, err)
			return
		}

		response.WithJSON(w, http.StatusOK, response.NewSuccessResponse())
	}
}

func HandleSearchAccounts(a *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sAccountTypeId := r.URL.Query().Get("accountTypeId")
		sExcludedAccountTypeID := r.URL.Query().Get("excludeAccountTypeId")

		session := auth.FromContext(r.Context())
		if session == nil {
			response.WithError(w, errors.ErrInvalidSession)
			return
		}
		filter := &account.AccountFilter{
			UserID: session.UserID,
		}

		var err error
		if sAccountTypeId != "" {
			filter.AccountTypeId, err = strconv.Atoi(sAccountTypeId)
			if err != nil {
				response.WithError(w, err)
				return
			}
		}

		if sExcludedAccountTypeID != "" {
			filter.ExcludedAccountTypeId, err = strconv.Atoi(sExcludedAccountTypeID)
			if err != nil {
				response.WithError(w, err)
				return
			}
		}

		results := a.AccountService.GetList(filter)
		response.WithJSON(w, http.StatusOK, results)
	}
}
