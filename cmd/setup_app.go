package main

import (
	"github.com/defryheryanto/piggy-bank-backend/internal/app"
)

func buildApp() *app.Application {
	piggyBankDatabase := setupDatabase()
	userConfigService := SetupUserConfigService(piggyBankDatabase)
	authService := SetupAuthService(piggyBankDatabase, userConfigService)
	accountService := setupAccountService(piggyBankDatabase)
	categoryService := SetupCategoryService(piggyBankDatabase)
	budgetService := SetupBudgetService(piggyBankDatabase, categoryService)
	transactionService := setupTransactionService(piggyBankDatabase)
	return &app.Application{
		AuthService:        authService,
		AccountService:     accountService,
		CategoryService:    categoryService,
		BudgetService:      budgetService,
		UserConfigService:  userConfigService,
		TransactionService: transactionService,
	}
}
