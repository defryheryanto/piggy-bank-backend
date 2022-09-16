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

func TestUpdateCategoryPayload(t *testing.T) {
	t.Run("Validate()", func(t *testing.T) {
		payload := &category.UpdateCategoryPayload{}

		err := payload.Validate()
		assert.NotNil(t, err)

		payload.CategoryId = 1
		err = payload.Validate()
		assert.NotNil(t, err)

		payload.CategoryName = "name"
		err = payload.Validate()
		assert.NotNil(t, err)

		payload.CategoryType = category.IncomeCategory
		err = payload.Validate()
		assert.Nil(t, err)
	})
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
			p := &category.UpdateCategoryPayload{
				CategoryId:   payload.CategoryId,
				CategoryName: "Salary",
				CategoryType: category.IncomeCategory,
				UserId:       int64(payload.UserId),
			}

			err := service.UpdateCategory(p)
			assert.NoError(t, err)

			var updated *category.Category
			db.Where("category_id = ?", p.CategoryId).First(&updated)
			assert.NotNil(t, updated)
			assert.Equal(t, p.CategoryName, updated.CategoryName)
			assert.Equal(t, p.CategoryType, updated.CategoryType)
			assert.Equal(t, p.UserId, int64(updated.UserId))
		})

		t.Run("should not update the budget", func(t *testing.T) {
			payloadWithBudget := &category.Category{
				CategoryId:   2,
				CategoryName: "Food",
				UserId:       1,
				CategoryType: category.ExpenseCategory,
				Budget:       500000,
			}

			db.Create(&payloadWithBudget)

			p := &category.UpdateCategoryPayload{
				CategoryId:   payloadWithBudget.CategoryId,
				CategoryName: "Salary",
				CategoryType: category.IncomeCategory,
				UserId:       int64(payload.UserId),
			}
			err := service.UpdateCategory(p)
			assert.NoError(t, err)

			var updated *category.Category
			db.Where("category_id = ?", p.CategoryId).First(&updated)
			assert.NotNil(t, updated)
			assert.Equal(t, p.CategoryName, updated.CategoryName)
			assert.Equal(t, p.CategoryType, updated.CategoryType)
			assert.Equal(t, payloadWithBudget.Budget, updated.Budget)
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

		t.Run("return error if category id is 0", func(t *testing.T) {
			p := &category.UpdateCategoryPayload{
				CategoryId:   0,
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
			p := &category.UpdateCategoryPayload{
				CategoryId:   payload.CategoryId,
				CategoryName: "Salary",
				CategoryType: "ini category",
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

func TestUpdateBudget(t *testing.T) {
	db := test.SetupTestDatabase(t, "../../.env", "../../db/migrations")

	test.RunTestWithDB(db, t, func(t *testing.T, db *gorm.DB) {
		service := setupService(db)

		payload := &category.Category{
			CategoryId:   1,
			CategoryName: "Food",
			UserId:       1,
			CategoryType: category.ExpenseCategory,
			Budget:       150000,
		}

		db.Create(&payload)

		t.Run("update the budget", func(t *testing.T) {
			err := service.UpdateBudget(payload.CategoryId, 200000)
			assert.NoError(t, err)

			updated, err := service.GetCategoryById(payload.CategoryId)
			assert.NoError(t, err)

			assert.Equal(t, int64(200000), updated.Budget)
		})

		t.Run("able to update the budget to 0", func(t *testing.T) {
			err := service.UpdateBudget(payload.CategoryId, 0)
			assert.NoError(t, err)

			updated, err := service.GetCategoryById(payload.CategoryId)
			assert.NoError(t, err)

			assert.Equal(t, int64(0), updated.Budget)
		})
	})
}

func TestGetList(t *testing.T) {
	db := test.SetupTestDatabase(t, "../../.env", "../../db/migrations")
	service := setupService(db)
	tables := []string{
		category.Category{}.TableName(),
	}

	test.TruncateAfterTest(t, db, tables, func() {
		expenseCategories := []*category.Category{
			{
				CategoryName: "Food",
				CategoryType: category.ExpenseCategory,
				UserId:       1,
				Budget:       500000,
			},
			{
				CategoryName: "Other",
				CategoryType: category.ExpenseCategory,
				UserId:       1,
				Budget:       2000000,
			},
			{
				CategoryName: "Insurance",
				CategoryType: category.ExpenseCategory,
				UserId:       2,
				Budget:       800000,
			},
		}
		db.Create(&expenseCategories)

		incomeCategories := []*category.Category{
			{
				CategoryName: "Salary",
				CategoryType: category.IncomeCategory,
				UserId:       1,
			},
		}
		db.Create(&incomeCategories)

		results := service.GetList(&category.CategoryFilter{
			UserId:       1,
			CategoryType: category.ExpenseCategory,
		})
		assert.Equal(t, 2, len(results))
	})
}
