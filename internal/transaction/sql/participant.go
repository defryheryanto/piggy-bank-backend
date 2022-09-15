package sql

import (
	"context"

	"github.com/defryheryanto/piggy-bank-backend/internal/storage"
	"github.com/defryheryanto/piggy-bank-backend/internal/transaction"
	"gorm.io/gorm"
)

type ParticipantStorage struct {
	db *gorm.DB
}

func NewParticipantStorage(db *gorm.DB) *ParticipantStorage {
	return &ParticipantStorage{db}
}

func (s *ParticipantStorage) BulkCreate(ctx context.Context, payloads []*transaction.Participant) error {
	db := s.getActiveDB(ctx)
	result := db.Create(&payloads)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *ParticipantStorage) getActiveDB(ctx context.Context) *gorm.DB {
	db := storage.DatabaseFromContext(ctx)
	if db == nil {
		return s.db
	}

	return db
}
