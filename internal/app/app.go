package app

import "github.com/defryheryanto/piggy-bank-backend/internal/auth"

type Application struct {
	AuthService *auth.AuthService
}
