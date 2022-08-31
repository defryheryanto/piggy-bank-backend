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

func TestGetCategoryById(t *testing.T) {
	db := test.SetupTestDatabase(t, "../../.env", "../../db/migrations")

	test.RunTestWithDB(db, t, func(t *testing.T, db *gorm.DB) {
		service := setupService(db)

		payload := &category.Category{
			CategoryId:   1,
			CategoryName: "Food",
			UserId:       1,
			CategoryType: category.ExpenseCategory,
		}

		db.Create(&payload)

		t.Run("should return category", func(t *testing.T) {
			cat, err := service.GetCategoryById(payload.CategoryId)
			assert.NoError(t, err)
			assert.NotNil(t, cat)
			assert.Equal(t, payload.CategoryId, cat.CategoryId)
		})

		t.Run("return error if category not found", func(t *testing.T) {
			cat, err := service.GetCategoryById(99)
			assert.ErrorIs(t, err, category.ErrCategoryNotFound)
			assert.Nil(t, cat)
		})
	})
}

func TestUpdateCategory(t *testing.T) {
	db := test.SetupTestDatabase(t, "../../.env", "../../db/migrations")

	test.RunTestWithDB(db, t, func(t *testing.T, db *gorm.DB) {
		service := setupService(db)

		payload := &category.Category{
			CategoryId:   1,
			CategoryName: "Food",
			UserId:       1,
			CategoryType: category.ExpenseCategory,
		}

		db.Create(&payload)

		t.Run("should update category", func(t *testing.T) {
			budget := int64(10000)
			p := &category.UpdateCategoryPayload{
				CategoryId:   payload.CategoryId,
				CategoryName: "Salary",
				CategoryType: category.IncomeCategory,
				Budget:       &budget,
				UserId:       int64(payload.UserId),
			}

			err := service.UpdateCategory(p)
			assert.NoError(t, err)

			var updated *category.Category
			db.Where("category_id = ?", p.CategoryId).First(&updated)
			assert.NotNil(t, updated)
			assert.Equal(t, p.CategoryName, updated.CategoryName)
			assert.Equal(t, p.CategoryType, updated.CategoryType)
			assert.Equal(t, *p.Budget, updated.Budget)
			assert.Equal(t, p.UserId, int64(updated.UserId))
		})

		t.Run("return error if category not found", func(t *testing.T) {
			p := &category.UpdateCategoryPayload{
				CategoryId:   99,
				CategoryName: "Salary",
				CategoryType: category.IncomeCategory,
				UserId:       int64(payload.UserId),
			}

			err := service.UpdateCategory(p)
			assert.ErrorIs(t, category.ErrCategoryNotFound, err)
		})

		t.Run("return error if user id not match with existing useri id", func(t *testing.T) {
			p := &category.UpdateCategoryPayload{
				CategoryId:   payload.CategoryId,
				CategoryName: "Salary",
				CategoryType: category.IncomeCategory,
				UserId:       99,
			}

			err := service.UpdateCategory(p)
			assert.ErrorIs(t, category.ErrCategoryNotFound, err)
		})

		t.Run("return error if category type is invalid", func(t *testing.T) {
			budget := int64(10000)
			p := &category.UpdateCategoryPayload{
				CategoryId:   payload.CategoryId,
				CategoryName: "Salary",
				CategoryType: "ini category",
				Budget:       &budget,
				UserId:       int64(payload.UserId),
			}

			err := service.UpdateCategory(p)
			assert.ErrorIs(t, category.ErrInvalidCategoryType, err)
		})
	})
}

func TestDeleteById(t *testing.T) {
	db := test.SetupTestDatabase(t, "../../.env", "../../db/migrations")

	test.RunTestWithDB(db, t, func(t *testing.T, db *gorm.DB) {
		service := setupService(db)

		payload := &category.Category{
			CategoryId:   1,
			CategoryName: "Food",
			UserId:       1,
			CategoryType: category.ExpenseCategory,
		}

		db.Create(&payload)

		t.Run("should delete existing data", func(t *testing.T) {
			err := service.DeleteById(payload.CategoryId, payload.UserId)
			assert.NoError(t, err)

			var existing *category.Category
			db.Where("category_id = ?", payload.CategoryId).First(&existing)
			assert.Equal(t, 0, existing.CategoryId)
		})

		t.Run("return error if category not found", func(t *testing.T) {
			err := service.DeleteById(99, payload.UserId)
			assert.ErrorIs(t, err, category.ErrCategoryNotFound)
		})

		t.Run("return error if user id is not match with existing", func(t *testing.T) {
			err := service.DeleteById(payload.CategoryId, 99)
			assert.ErrorIs(t, err, category.ErrCategoryNotFound)
		})
	})
}
