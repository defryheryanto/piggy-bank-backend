package sql

import (
	"github.com/defryheryanto/piggy-bank-backend/internal/transaction"
	"gorm.io/gorm"
)

type ParticipantStorage struct {
	db *gorm.DB
}

func NewParticipantStorage(db *gorm.DB) *ParticipantStorage {
	return &ParticipantStorage{db}
}

func (s *ParticipantStorage) BulkCreate(payloads []*transaction.Participant) error {
	result := s.db.Create(&payloads)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
