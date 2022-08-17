package main

import (
	"github.com/defryheryanto/piggy-bank-backend/internal/app"
)

func buildApp() *app.Application {
	piggyBankDatabase := setupDatabase()
	authService := SetupAuthService(piggyBankDatabase)
	return &app.Application{
		AuthService: authService,
	}
}
