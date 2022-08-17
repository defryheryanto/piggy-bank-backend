package handler

import (
	"net/http"

	"github.com/defryheryanto/piggy-bank-backend/internal/app"
	"github.com/defryheryanto/piggy-bank-backend/internal/httpserver/response"
)

func HandleGetTypes(a *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		types := a.AccountService.GetTypes()
		response.WithJSON(w, http.StatusOK, types)
	}
}
