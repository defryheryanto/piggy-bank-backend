package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/defryheryanto/piggy-bank-backend/internal/app"
	"github.com/defryheryanto/piggy-bank-backend/internal/auth"
	"github.com/defryheryanto/piggy-bank-backend/internal/budget"
	"github.com/defryheryanto/piggy-bank-backend/internal/errors"
	"github.com/defryheryanto/piggy-bank-backend/internal/httpserver/response"
)

func HandleCreateBudget(a *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := &budget.CreateBudgetPayload{}

		err := json.NewDecoder(r.Body).Decode(payload)
		if err != nil {
			if err == io.EOF {
				response.WithError(w, errors.ErrEmptyPayload)
				return
			}
			response.WithError(w, errors.ErrUnprocessablePayload)
			return
		}

		category, err := a.CategoryService.GetCategoryById(payload.CategoryId)
		if err != nil {
			response.WithError(w, err)
			return
		}

		session := auth.FromContext(r.Context())
		if session == nil {
			response.WithError(w, errors.ErrInvalidSession)
			return
		}
		if session.UserID != category.UserId {
			response.WithError(w, errors.NewForbiddenError("Forbidden", "unable to create budget with other user's category"))
			return
		}

		err = a.BudgetService.Create(payload)
		if err != nil {
			response.WithError(w, err)
			return
		}

		response.WithJSON(w, http.StatusOK, response.NewSuccessResponse())
	}
}
