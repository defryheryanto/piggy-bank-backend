package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/defryheryanto/piggy-bank-backend/internal/app"
	"github.com/defryheryanto/piggy-bank-backend/internal/auth"
	"github.com/defryheryanto/piggy-bank-backend/internal/errors"
	"github.com/defryheryanto/piggy-bank-backend/internal/httpserver/response"
	"github.com/defryheryanto/piggy-bank-backend/internal/transaction"
)

// Basic Transaction is transaction with type income or expense
func HandleCreateBasicTransaction(a *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload *transaction.CreateBasicTransactionPayload

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
		payload.TransactionDate = time.Now()

		err = a.TransactionService.CreateBasic(r.Context(), payload)
		if err != nil {
			response.WithError(w, err)
			return
		}

		response.WithJSON(w, http.StatusOK, response.NewSuccessResponse())
	}
}

func HandleCreateTransfer(a *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload *transaction.CreateTransferPayload

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
		payload.TransactionDate = time.Now()

		err = a.TransactionService.CreateTransfer(r.Context(), payload)
		if err != nil {
			response.WithError(w, err)
			return
		}

		response.WithJSON(w, http.StatusOK, response.NewSuccessResponse())
	}
}

func HandleCreateSavingTransaction(a *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload *transaction.CreateSavingPayload

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
		payload.TransactionDate = time.Now()

		err = a.TransactionService.CreateSaving(r.Context(), payload)
		if err != nil {
			response.WithError(w, err)
			return
		}

		response.WithJSON(w, http.StatusOK, response.NewSuccessResponse())
	}
}
