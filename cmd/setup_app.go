package main

import (
	"github.com/defryheryanto/piggy-bank-backend/internal/app"
)

func buildApp() *app.Application {
	piggyBankDatabase := setupDatabase()
	authService := SetupAuthService(piggyBankDatabase)
	accountService := setupAccountService(piggyBankDatabase)
	categoryService := SetupCategoryService(piggyBankDatabase)
	budgetService := SetupBudgetService(piggyBankDatabase, categoryService)
	return &app.Application{
		AuthService:     authService,
		AccountService:  accountService,
		CategoryService: categoryService,
		BudgetService:   budgetService,
	}
}
