package httpserver

import (
	"net/http"

	"github.com/defryheryanto/piggy-bank-backend/internal/account/handler"
	auth_handler "github.com/defryheryanto/piggy-bank-backend/internal/auth/handler"
	"github.com/defryheryanto/piggy-bank-backend/internal/httpserver/middleware"
	"github.com/defryheryanto/piggy-bank-backend/internal/httpserver/response"
	"github.com/gorilla/mux"
)

func (s *ApplicationServer) CompileRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		response.WithJSON(rw, 200, map[string]string{
			"status": "healhty",
		})
	})

	r.HandleFunc("/api/v1/register", auth_handler.HandleRegister(s.application)).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/login", auth_handler.HandleLogin(s.application)).Methods(http.MethodPost)

	privateRoute := r.NewRoute().Subrouter()
	privateRoute.Use(middleware.PrivateRoute(s.application))
	privateRoute.HandleFunc("/api/v1/accounts/types", handler.HandleGetTypes(s.application)).Methods(http.MethodGet)

	return r
}
