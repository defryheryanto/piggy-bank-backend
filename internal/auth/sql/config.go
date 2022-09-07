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

func (s *UserConfigStorage) GetByUserId(userId int) *auth.UserConfig {
	var cfg *auth.UserConfig

	s.db.Where("user_id = ?", userId).First(&cfg)
	if cfg.ConfigId == 0 {
		return nil
	}

	return cfg
}

func (s *UserConfigStorage) Update(payload *auth.UserConfig) error {
	db := s.db.Where("user_id = ?", payload.UserId).Updates(&payload)
	if db.Error != nil {
		return db.Error
	}

	return nil
}
