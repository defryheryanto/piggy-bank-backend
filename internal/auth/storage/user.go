package storage

import (
	"github.com/defryheryanto/piggy-bank-backend/internal/auth"
	"gorm.io/gorm"
)

type UserStorage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) *UserStorage {
	return &UserStorage{db}
}

func (s *UserStorage) Create(payload *auth.User) error {
	db := s.db.Create(&payload)
	if db.Error != nil {
		return db.Error
	}

	return nil
}

func (s *UserStorage) GetByUsername(username string) *auth.User {
	var data *auth.User

	s.db.Where("username = ?", username).Find(&data)
	if data.UserID == 0 {
		return nil
	}

	return data
}
