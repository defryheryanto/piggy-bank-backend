package sql

import (
	"github.com/defryheryanto/piggy-bank-backend/internal/transaction"
	"gorm.io/gorm"
)

type TransactionStorage struct {
	db *gorm.DB
}

func NewTransactionStorage(db *gorm.DB) *TransactionStorage {
	return &TransactionStorage{db}
}

func (s *TransactionStorage) Create(payload *transaction.Transaction) error {
	result := s.db.Create(&payload)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
