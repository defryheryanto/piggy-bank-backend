package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/defryheryanto/piggy-bank-backend/internal/app"
	"github.com/defryheryanto/piggy-bank-backend/internal/auth"
	"github.com/defryheryanto/piggy-bank-backend/internal/category"
	"github.com/defryheryanto/piggy-bank-backend/internal/errors"
	"github.com/defryheryanto/piggy-bank-backend/internal/httpserver/response"
	"github.com/gorilla/mux"
)

func HandleCreateCategory(a *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p *CreateCategoryPayload

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			if err == io.EOF {
				response.WithError(w, errors.ErrEmptyPayload)
				return
			}
			response.WithError(w, errors.ErrUnprocessablePayload)
			return
		}

		err = p.Validate()
		if err != nil {
			response.WithError(w, err)
			return
		}

		session := auth.FromContext(r.Context())
		if session == nil {
			response.WithError(w, errors.ErrInvalidSession)
			return
		}

		payload := &category.Category{
			CategoryName: p.CategoryName,
			CategoryType: category.CategoryType(p.CategoryType),
			UserId:       session.UserID,
			Budget:       p.Budget,
		}

		err = a.CategoryService.Create(payload)
		if err != nil {
			response.WithError(w, err)
			return
		}

		response.WithJSON(w, http.StatusOK, response.NewSuccessResponse())
	}
}

func HandleGetCategoryTypes(a *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := auth.FromContext(r.Context())
		if session == nil {
			response.WithError(w, errors.ErrInvalidSession)
			return
		}

		result := a.CategoryService.GetCategoryTypeDetails(session.UserID)

		response.WithJSON(w, http.StatusOK, result)
	}
}

func HandleGetCategory(a *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sCategoryId := mux.Vars(r)["category_id"]
		categoryId, err := strconv.Atoi(sCategoryId)
		if err != nil {
			response.WithError(w, errors.NewBadRequestError("Category ID must be an number", "category_id is not a number"))
			return
		}

		category, err := a.CategoryService.GetCategoryById(categoryId)
		if err != nil {
			response.WithError(w, err)
			return
		}

		response.WithJSON(w, http.StatusOK, category)
	}
}

func HandleUpdateCategory(a *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sCategoryId := mux.Vars(r)["category_id"]
		categoryId, err := strconv.Atoi(sCategoryId)
		if err != nil {
			response.WithError(w, errors.NewBadRequestError("Category ID must be an number", "category_id is not a number"))
			return
		}

		var p *category.UpdateCategoryPayload
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
		p.CategoryId = categoryId
		p.UserId = int64(session.UserID)
		err = p.Validate()
		if err != nil {
			response.WithError(w, err)
			return
		}

		err = a.CategoryService.UpdateCategory(p)
		if err != nil {
			response.WithError(w, err)
			return
		}

		response.WithJSON(w, http.StatusOK, response.NewSuccessResponse())
	}
}

func HandleDeleteCategory(a *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sCategoryId := mux.Vars(r)["category_id"]
		categoryId, err := strconv.Atoi(sCategoryId)
		if err != nil {
			response.WithError(w, errors.NewBadRequestError("Category ID must be an number", "category_id is not a number"))
			return
		}

		session := auth.FromContext(r.Context())
		if session == nil {
			response.WithError(w, errors.ErrInvalidSession)
			return
		}

		err = a.CategoryService.DeleteById(categoryId, session.UserID)
		if err != nil {
			response.WithError(w, err)
			return
		}

		response.WithJSON(w, http.StatusOK, response.NewSuccessResponse())
	}
}
