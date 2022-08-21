package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/defryheryanto/piggy-bank-backend/internal/app"
	"github.com/defryheryanto/piggy-bank-backend/internal/auth"
	"github.com/defryheryanto/piggy-bank-backend/internal/category"
	"github.com/defryheryanto/piggy-bank-backend/internal/errors"
	"github.com/defryheryanto/piggy-bank-backend/internal/httpserver/response"
)

func HandleCreateCategory(a *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p *CreateCategoryPayload

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			if err == io.EOF {
				response.WithError(w, errors.EmptyPayload)
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
			response.WithError(w, errors.InvalidSession)
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
