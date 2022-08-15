package httpserver

import (
	"net/http"

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

	return r
}
