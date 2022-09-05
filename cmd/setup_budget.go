package main

import (
	"github.com/defryheryanto/piggy-bank-backend/internal/budget"
	budget_storage "github.com/defryheryanto/piggy-bank-backend/internal/budget/sql"
	"github.com/defryheryanto/piggy-bank-backend/internal/category"
	"gorm.io/gorm"
)

func SetupBudgetService(db *gorm.DB, category *category.CategoryService) *budget.BudgetService {
	budgetStorage := budget_storage.NewBudgetStorage(db)

	return budget.NewBudgetService(budgetStorage, category)
}
