package sql

import (
	"context"

	"github.com/defryheryanto/piggy-bank-backend/internal/transaction"
)

type TransactionMockStorage struct {
	mock MockFunction
}

func NewTransactionMockStorage(fn MockFunction) *TransactionMockStorage {
	return &TransactionMockStorage{fn}
}

func (s *TransactionMockStorage) Create(ctx context.Context, payload *transaction.Transaction) error {
	return s.mock()
}
