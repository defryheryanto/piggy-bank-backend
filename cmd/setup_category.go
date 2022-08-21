package main

import (
	"github.com/defryheryanto/piggy-bank-backend/internal/category"
	category_storage "github.com/defryheryanto/piggy-bank-backend/internal/category/sql"
	"gorm.io/gorm"
)

func SetupCategoryService(db *gorm.DB) *category.CategoryService {
	storage := category_storage.NewCategoryStorage(db)

	return category.NewCategoryService(storage)
}
