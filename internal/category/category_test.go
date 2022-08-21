package category_test

import (
	"testing"

	"github.com/defryheryanto/piggy-bank-backend/internal/category"
	category_storage "github.com/defryheryanto/piggy-bank-backend/internal/category/sql"
	"github.com/defryheryanto/piggy-bank-backend/test"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupService(db *gorm.DB) *category.CategoryService {
	userStorage := category_storage.NewCategoryStorage(db)

	return category.NewCategoryService(userStorage)
}

func TestCreate(t *testing.T) {
	db := test.SetupTestDatabase(t, "../../.env", "../../db/migrations")

	test.RunTestWithDB(db, t, func(t *testing.T, db *gorm.DB) {
		service := setupService(db)

		t.Run("should insert to db", func(t *testing.T) {
			payload := &category.Category{
				CategoryName: "Food",
				UserId:       1,
				CategoryType: category.ExpenseCategory,
			}

			err := service.Create(payload)
			assert.NoError(t, err)

			var data *category.Category
			db.Where("category_name = ?", payload.CategoryName).First(&data)

			if data.CategoryId == 0 {
				t.Errorf("category not inserted to database")
			}

			assert.Equal(t, payload.CategoryType, data.CategoryType)
			assert.Equal(t, payload.UserId, data.UserId)
			assert.Equal(t, int64(0), data.Budget)
		})

		t.Run("return error if type is invalid", func(t *testing.T) {
			payload := &category.Category{
				CategoryName: "Food",
				UserId:       1,
				CategoryType: "pemasukan",
			}

			err := service.Create(payload)
			assert.ErrorIs(t, err, category.ErrInvalidCategoryType)
		})
	})
}

func TestGetTypeDetails(t *testing.T) {
	db := test.SetupTestDatabase(t, "../../.env", "../../db/migrations")

	test.RunTestWithDB(db, t, func(t *testing.T, db *gorm.DB) {
		service := setupService(db)

		payload := []*category.Category{
			{
				CategoryName: "Salary",
				UserId:       1,
				CategoryType: category.IncomeCategory,
				Budget:       1000000,
			},
			{
				CategoryName: "Food",
				UserId:       1,
				CategoryType: category.ExpenseCategory,
				Budget:       1000000,
			},
			{
				CategoryName: "Food",
				UserId:       2,
				CategoryType: category.ExpenseCategory,
			},
		}

		db.Create(&payload)

		t.Run("should get by available type", func(t *testing.T) {
			result := service.GetCategoryTypeDetails(1)
			assert.Equal(t, len(category.CategoryTypes), len(result))
		})

		t.Run("should get categories", func(t *testing.T) {
			result := service.GetCategoryTypeDetails(1)

			//check income categories
			assert.Equal(t, 1, len(result[0].Categories))

			//check expense categories
			assert.Equal(t, 1, len(result[1].Categories))
		})
	})
}
