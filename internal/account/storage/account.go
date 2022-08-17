package storage

import (
	"github.com/defryheryanto/piggy-bank-backend/internal/account"
	"gorm.io/gorm"
)

type AccountStorage struct {
	db *gorm.DB
}

func NewAccountStorage(db *gorm.DB) *AccountStorage {
	return &AccountStorage{db}
}

func (s *AccountStorage) GetTypes() []*account.AccountType {
	var datas []*account.AccountType
	s.db.Find(&datas)
	return datas
}
