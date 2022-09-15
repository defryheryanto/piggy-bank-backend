package sql

import (
	"context"

	"github.com/defryheryanto/piggy-bank-backend/internal/transaction"
)

type ParticipantMockStorage struct {
	mock MockFunction
}

func NewParticipantMockStorage(fn MockFunction) *ParticipantMockStorage {
	return &ParticipantMockStorage{fn}
}

func (s *ParticipantMockStorage) BulkCreate(ctx context.Context, payloads []*transaction.Participant) error {
	return s.mock()
}
