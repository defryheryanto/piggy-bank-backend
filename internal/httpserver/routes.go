package httpserver

import (
	"net/http"

	account_handler "github.com/defryheryanto/piggy-bank-backend/internal/account/handler"
	auth_handler "github.com/defryheryanto/piggy-bank-backend/internal/auth/handler"
	budget_handler "github.com/defryheryanto/piggy-bank-backend/internal/budget/handler"
	category_handler "github.com/defryheryanto/piggy-bank-backend/internal/category/handler"
	"github.com/defryheryanto/piggy-bank-backend/internal/httpserver/middleware"
	"github.com/defryheryanto/piggy-bank-backend/internal/httpserver/response"
	"github.com/gorilla/mux"
)

func (s *ApplicationServer) CompileRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		response.WithJSON(rw, 200, map[string]string{
			"status": "healthy",
		})
	})

	r.HandleFunc("/api/v1/register", auth_handler.HandleRegister(s.application)).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/login", auth_handler.HandleLogin(s.application)).Methods(http.MethodPost)

	privateRoute := r.NewRoute().Subrouter()
	privateRoute.Use(middleware.PrivateRoute(s.application))
	privateRoute.HandleFunc("/api/v1/accounts", account_handler.HandleGetAccounts(s.application)).Methods(http.MethodGet)
	privateRoute.HandleFunc("/api/v1/accounts", account_handler.HandleCreateAccount(s.application)).Methods(http.MethodPost)
	privateRoute.HandleFunc("/api/v1/accounts/{account_id}", account_handler.HandleUpdateAccount(s.application)).Methods(http.MethodPatch)
	privateRoute.HandleFunc("/api/v1/accounts/{account_id}", account_handler.HandleDeleteAccount(s.application)).Methods(http.MethodDelete)
	privateRoute.HandleFunc("/api/v1/accounts/types", account_handler.HandleGetTypes(s.application)).Methods(http.MethodGet)

	privateRoute.HandleFunc("/api/v1/categories", category_handler.HandleCreateCategory(s.application)).Methods(http.MethodPost)
	privateRoute.HandleFunc("/api/v1/categories", category_handler.HandleGetCategoryTypes(s.application)).Methods(http.MethodGet)
	privateRoute.HandleFunc("/api/v1/categories/{category_id}", category_handler.HandleGetCategory(s.application)).Methods(http.MethodGet)
	privateRoute.HandleFunc("/api/v1/categories/{category_id}", category_handler.HandleUpdateCategory(s.application)).Methods(http.MethodPatch)
	privateRoute.HandleFunc("/api/v1/categories/{category_id}", category_handler.HandleDeleteCategory(s.application)).Methods(http.MethodDelete)

	privateRoute.HandleFunc("/api/v1/budgets", budget_handler.HandleUpsertBudget(s.application)).Methods(http.MethodPut)
	privateRoute.HandleFunc("/api/v1/categories/{category_id}/budgets", budget_handler.HandleGetBudgetYearSummary(s.application)).Methods(http.MethodGet)

	return r
}
