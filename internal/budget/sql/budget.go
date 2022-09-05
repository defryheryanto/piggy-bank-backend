package sql

import (
	"github.com/defryheryanto/piggy-bank-backend/internal/budget"
	"gorm.io/gorm"
)

type BudgetStorage struct {
	db *gorm.DB
}

func NewBudgetStorage(db *gorm.DB) *BudgetStorage {
	return &BudgetStorage{db}
}

func (s *BudgetStorage) Create(payload *budget.Budget) error {
	result := s.db.Create(&payload)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *BudgetStorage) Update(payload *budget.Budget) error {
	result := s.db.Where("budget_id = ?", payload.BudgetId).Updates(&payload)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *BudgetStorage) GetByMonthAndYear(categoryId, month, year int) *budget.Budget {
	var budget *budget.Budget

	s.db.Where("month = ? AND year = ? AND category_id = ?", month, year, categoryId).First(&budget)
	if budget.BudgetId == 0 {
		return nil
	}

	return budget
}
