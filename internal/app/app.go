package app

import (
	"github.com/defryheryanto/piggy-bank-backend/internal/account"
	"github.com/defryheryanto/piggy-bank-backend/internal/auth"
)

type Application struct {
	AuthService    *auth.AuthService
	AccountService *account.AccountService
}
