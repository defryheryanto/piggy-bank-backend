package main

import (
	"github.com/defryheryanto/piggy-bank-backend/internal/budget"
	budget_storage "github.com/defryheryanto/piggy-bank-backend/internal/budget/sql"
	"gorm.io/gorm"
)

func SetupBudgetService(db *gorm.DB) *budget.BudgetService {
	budgetStorage := budget_storage.NewBudgetStorage(db)

	return budget.NewBudgetService(budgetStorage)
}
