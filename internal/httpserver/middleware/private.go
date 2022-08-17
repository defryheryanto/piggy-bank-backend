package middleware

import (
	"net/http"
	"strings"

	"github.com/defryheryanto/piggy-bank-backend/internal/app"
	"github.com/defryheryanto/piggy-bank-backend/internal/errors"
	"github.com/defryheryanto/piggy-bank-backend/internal/httpserver/response"
	"github.com/gorilla/mux"
)

func PrivateRoute(a *app.Application) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			tokens := strings.Split(authorization, "Bearer ")
			if len(tokens) < 2 {
				response.WithError(w, errors.NewUnauthorizedError("Unauthorized", "invalid token"))
				return
			}

			token := tokens[1]
			isValid, err := a.AuthService.Authenticate(token)
			if err != nil {
				response.WithError(w, errors.NewUnauthorizedError("Unauthorized", err.Error()))
				return
			}

			if !isValid {
				response.WithError(w, errors.NewUnauthorizedError("Your session has expired, Please Re-Login", "invalid token"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
