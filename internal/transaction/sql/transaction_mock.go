package sql

import (
	"context"

	"github.com/defryheryanto/piggy-bank-backend/internal/transaction"
)

type TransactionMockStorage struct {
	create               MockFunction
	createTransferDetail MockFunction
}

func NewTransactionMockStorage(create, createTransferDetail MockFunction) *TransactionMockStorage {
	return &TransactionMockStorage{create, createTransferDetail}
}

func (s *TransactionMockStorage) Create(ctx context.Context, payload *transaction.Transaction) error {
	return s.create()
}

func (s *TransactionMockStorage) CreateTransferDetail(ctx context.Context, payload *transaction.TransferDetail) error {
	return s.createTransferDetail()
}
