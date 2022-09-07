package sql

import (
	"github.com/defryheryanto/piggy-bank-backend/internal/auth"
	"gorm.io/gorm"
)

type UserConfigStorage struct {
	db *gorm.DB
}

func NewUserConfigStorage(db *gorm.DB) *UserConfigStorage {
	return &UserConfigStorage{db}
}

func (s *UserConfigStorage) Create(payload *auth.UserConfig) error {
	db := s.db.Create(&payload)
	if db.Error != nil {
		return db.Error
	}

	return nil
}
