package sql

import (
	"strings"

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

func (s *CategoryStorage) GetByTypeAndUserId(categoryType category.CategoryType, userId int) []*category.Category {
	var datas []*category.Category
	s.db.Where("category_type = ? AND user_id = ?", categoryType, userId).Find(&datas)

	return datas
}

func (s *CategoryStorage) GetById(categoryId int) *category.Category {
	var data *category.Category

	s.db.Where("category_id = ?", categoryId).Find(&data)
	if data.CategoryId == 0 {
		return nil
	}

	return data
}

func (s *CategoryStorage) UpdateById(categoryId int, payload *category.Category) error {
	res := s.db.Model(payload).
		Select("category_name", "category_type", "budget").
		Where("category_id = ?", categoryId).
		Updates(map[string]interface{}{
			"category_name": payload.CategoryName,
			"category_type": payload.CategoryType,
			"budget":        payload.Budget,
		})
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (s *CategoryStorage) DeleteById(categoryId int) error {
	err := s.db.Model(&category.Category{}).Delete("category_id = ?", categoryId)
	if err.Error != nil {
		return err.Error
	}

	return nil
}

func (s *CategoryStorage) GetByFilter(filter *category.CategoryFilter) []*category.Category {
	result := []*category.Category{}
	if filter == nil {
		s.db.Find(&result)
		return result
	}

	whereQueries := []string{}
	values := []interface{}{}

	if filter.CategoryType != "" {
		whereQueries = append(whereQueries, "category_type = ?")
		values = append(values, string(filter.CategoryType))
	}
	if filter.UserId != 0 {
		whereQueries = append(whereQueries, "user_id = ?")
		values = append(values, filter.UserId)
	}

	s.db.Where(strings.Join(whereQueries, " AND "), values...).Find(&result)
	return result
}
