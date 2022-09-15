package sql

import (
	"context"

	"github.com/defryheryanto/piggy-bank-backend/internal/storage"
	"github.com/defryheryanto/piggy-bank-backend/internal/transaction"
	"gorm.io/gorm"
)

type TransactionStorage struct {
	db *gorm.DB
}

func NewTransactionStorage(db *gorm.DB) *TransactionStorage {
	return &TransactionStorage{db}
}

func (s *TransactionStorage) Create(ctx context.Context, payload *transaction.Transaction) error {
	db := s.getActiveDB(ctx)
	result := db.Create(&payload)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *TransactionStorage) getActiveDB(ctx context.Context) *gorm.DB {
	db := storage.DatabaseFromContext(ctx)
	if db == nil {
		return s.db
	}

	return db
}
