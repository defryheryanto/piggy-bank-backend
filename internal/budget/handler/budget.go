package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/defryheryanto/piggy-bank-backend/internal/app"
	"github.com/defryheryanto/piggy-bank-backend/internal/auth"
	"github.com/defryheryanto/piggy-bank-backend/internal/budget"
	"github.com/defryheryanto/piggy-bank-backend/internal/errors"
	"github.com/defryheryanto/piggy-bank-backend/internal/httpserver/response"
	"github.com/gorilla/mux"
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

func HandleGetBudgetYearSummary(a *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sCategoryId := mux.Vars(r)["category_id"]
		sYear := r.URL.Query().Get("year")

		categoryId, err := strconv.Atoi(sCategoryId)
		if err != nil {
			response.WithError(w, errors.NewBadRequestError("Category Id must be a number", "category_id must be a numeric"))
			return
		}

		year, err := strconv.Atoi(sYear)
		if err != nil {
			response.WithError(w, errors.NewBadRequestError("Year must be a number", "year must be a numeric"))
			return
		}

		summary, err := a.BudgetService.GetBudgetYearSummary(categoryId, year)
		if err != nil {
			response.WithError(w, err)
			return
		}

		response.WithJSON(w, http.StatusOK, summary)
	}
}
