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
