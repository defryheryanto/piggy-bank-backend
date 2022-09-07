package app

import (
	"github.com/defryheryanto/piggy-bank-backend/internal/account"
	"github.com/defryheryanto/piggy-bank-backend/internal/auth"
	"github.com/defryheryanto/piggy-bank-backend/internal/budget"
	"github.com/defryheryanto/piggy-bank-backend/internal/category"
)

type Application struct {
	AuthService       *auth.AuthService
	AccountService    *account.AccountService
	CategoryService   *category.CategoryService
	BudgetService     *budget.BudgetService
	UserConfigService *auth.UserConfigService
}
