package sql

import (
	"github.com/defryheryanto/piggy-bank-backend/internal/category"
	"gorm.io/gorm"
)

type CategoryStorage struct {
	db *gorm.DB
}

func NewCategoryStorage(db *gorm.DB) *CategoryStorage {
	return &CategoryStorage{db}
}

func (s *CategoryStorage) Create(payload *category.Category) error {
	res := s.db.Create(&payload)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
