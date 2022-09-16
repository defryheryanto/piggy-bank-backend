package sql

import (
	"context"

	"github.com/defryheryanto/piggy-bank-backend/internal/transaction"
)

type TransactionMockStorage struct {
	createFunc               MockFunction
	createTransferDetailFunc MockFunction
	createSavingDetailFunc   MockFunction
}

func NewTransactionMockStorage(create, createTransferDetail, createSavingDetail MockFunction) *TransactionMockStorage {
	return &TransactionMockStorage{create, createTransferDetail, createSavingDetail}
}

func (s *TransactionMockStorage) Create(ctx context.Context, payload *transaction.Transaction) error {
	return s.createFunc()
}

func (s *TransactionMockStorage) CreateTransferDetail(ctx context.Context, payload *transaction.TransferDetail) error {
	return s.createTransferDetailFunc()
}

func (s *TransactionMockStorage) CreateSavingDetail(ctx context.Context, payload *transaction.SavingDetail) error {
	return s.createSavingDetailFunc()
}
