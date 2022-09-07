package budget_test

import (
	"testing"
	"time"

	"github.com/defryheryanto/piggy-bank-backend/internal/budget"
	budget_storage "github.com/defryheryanto/piggy-bank-backend/internal/budget/sql"
	"github.com/defryheryanto/piggy-bank-backend/internal/category"
	category_storage "github.com/defryheryanto/piggy-bank-backend/internal/category/sql"
	"github.com/defryheryanto/piggy-bank-backend/test"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

func setupService(db *gorm.DB) *budget.BudgetService {
	budgetStorage := budget_storage.NewBudgetStorage(db)
	categoryStorage := category_storage.NewCategoryStorage(db)
	category := category.NewCategoryService(categoryStorage)

	return budget.NewBudgetService(budgetStorage, category)
}

func setupDatabase(t *testing.T) *gorm.DB {
	db := test.SetupTestDatabase(t, "../../.env", "../../db/migrations")
	return db
}

func initBudget(db *gorm.DB, t *testing.T) (*category.Category, []*budget.Budget) {
	initCategory := &category.Category{
		CategoryId:   1,
		CategoryName: "Food",
		CategoryType: category.ExpenseCategory,
		Budget:       2000000,
		UserId:       1,
	}
	res := db.Create(&initCategory)
	assert.NoError(t, res.Error)

	initBudgets := []*budget.Budget{
		{
			CategoryId: initCategory.CategoryId,
			Month:      1,
			Year:       2022,
			Budget:     1500000,
		},
		{
			CategoryId: initCategory.CategoryId,
			Month:      2,
			Year:       2022,
			Budget:     2100000,
		},
		{
			CategoryId: initCategory.CategoryId,
			Month:      5,
			Year:       2022,
			Budget:     500000,
		},
		{
			CategoryId: initCategory.CategoryId,
			Month:      12,
			Year:       2022,
			Budget:     4500000,
		},
	}
	res = db.Create(&initBudgets)
	assert.NoError(t, res.Error)

	return initCategory, initBudgets
}

func TestCreateOrUpdate(t *testing.T) {
	db := setupDatabase(t)

	test.RunTestWithDB(db, t, func(t *testing.T, db *gorm.DB) {
		service := setupService(db)

		t.Run("should insert to db if not exists", func(t *testing.T) {
			payload := &budget.CreateBudgetPayload{
				CategoryId: 1,
				Month:      9,
				Year:       2022,
				Budget:     1500000,
			}

			err := service.CreateOrUpdate(payload)
			assert.NoError(t, err)
		})

		t.Run("should update db if existing month and year budget is exists", func(t *testing.T) {
			initBudget := &budget.Budget{
				BudgetId:   77,
				CategoryId: 1,
				Month:      1,
				Year:       2022,
				Budget:     4000000,
			}

			res := db.Create(&initBudget)
			assert.NoError(t, res.Error)

			payload := &budget.CreateBudgetPayload{
				CategoryId: initBudget.CategoryId,
				Month:      initBudget.Month,
				Year:       initBudget.Year,
				Budget:     1245000,
			}
			err := service.CreateOrUpdate(payload)
			assert.NoError(t, err)

			existingBudget := &budget.Budget{}
			db.Where("budget_id = ?", initBudget.BudgetId).First(&existingBudget)
			assert.Equal(t, payload.Budget, existingBudget.Budget)
		})

		t.Run("return error if month is invalid", func(t *testing.T) {
			payload := &budget.CreateBudgetPayload{
				CategoryId: 2,
				Month:      0,
				Year:       2022,
				Budget:     1500000,
			}

			err := service.CreateOrUpdate(payload)
			assert.ErrorIs(t, err, budget.ErrInvalidMonthBudget)

			payload = &budget.CreateBudgetPayload{
				CategoryId: 2,
				Month:      13,
				Year:       2022,
				Budget:     1500000,
			}
			err = service.CreateOrUpdate(payload)
			assert.ErrorIs(t, err, budget.ErrInvalidMonthBudget)
		})
	})
}

func TestGetBudgetYearSummary(t *testing.T) {
	db := setupDatabase(t)
	test.RunTestWithDB(db, t, func(t *testing.T, db *gorm.DB) {
		service := setupService(db)

		initCategory, initBudgets := initBudget(db, t)

		t.Run("should get budget summary defaulted to category's budget", func(t *testing.T) {
			summary, err := service.GetBudgetYearSummary(1, 2022)
			assert.NoError(t, err)

			assert.Equal(t, initCategory.Budget, summary.DefaultBudget)
			assert.Equal(t, 12, len(summary.Budgets))

			for _, detail := range summary.Budgets {
				index := slices.IndexFunc(initBudgets, func(b *budget.Budget) bool {
					return time.Month(b.Month).String() == detail.Month
				})
				if index == -1 {
					assert.Equal(t, initCategory.Budget, detail.Budget)
				} else {
					assert.Equal(t, initBudgets[index].Budget, detail.Budget)
				}
			}
		})

		t.Run("return error if category not found", func(t *testing.T) {
			summary, err := service.GetBudgetYearSummary(99, 2022)
			assert.ErrorIs(t, err, category.ErrCategoryNotFound)
			assert.Nil(t, summary)
		})
	})
}

func TestGetActiveBudge(t *testing.T) {
	db := setupDatabase(t)

	test.RunTestWithDB(db, t, func(t *testing.T, db *gorm.DB) {
		service := setupService(db)

		initCategories := []*category.Category{
			{
				CategoryId:   1,
				CategoryName: "Food",
				UserId:       1,
				CategoryType: category.ExpenseCategory,
				Budget:       1500000,
			},
			{
				CategoryId:   2,
				CategoryName: "Food",
				UserId:       2,
				CategoryType: category.ExpenseCategory,
			},
			{
				CategoryId:   3,
				CategoryName: "Other",
				UserId:       1,
				CategoryType: category.ExpenseCategory,
				Budget:       2500000,
			},
			{
				CategoryId:   4,
				CategoryName: "Salary",
				UserId:       1,
				CategoryType: category.IncomeCategory,
				Budget:       25000000,
			},
		}

		db.Create(&initCategories)

		t.Run("get active budget on expense category only", func(t *testing.T) {
			actives := service.GetActiveBudget(1, 9, 2022)

			//check if length of active budget is the same with expense categories count
			assert.Equal(t, 2, len(actives))
		})

		t.Run("get active budget from category if budget not exists", func(t *testing.T) {
			actives := service.GetActiveBudget(1, 9, 2022)

			//check if length of active budget is the same with expense categories count
			assert.Equal(t, 2, len(actives))

			for _, active := range actives {
				for _, cat := range initCategories {
					if active.CategoryName == cat.CategoryName && cat.UserId == 1 {
						assert.Equal(t, cat.Budget, active.Budget)
					}
				}
			}
		})

		t.Run("get active budget from budget if exists", func(t *testing.T) {
			budgets := []*budget.Budget{
				{
					CategoryId: 1,
					Month:      9,
					Year:       2022,
					Budget:     9500000,
				},
				{
					CategoryId: 3,
					Month:      9,
					Year:       2022,
					Budget:     2150000,
				},
			}

			db.Create(budgets)
			actives := service.GetActiveBudget(1, 9, 2022)

			//check if length of active budget is the same with expense categories count
			assert.Equal(t, 2, len(actives))

			for _, active := range actives {
				for _, bdg := range budgets {
					if active.CategoryId == bdg.CategoryId {
						assert.Equal(t, bdg.Budget, active.Budget)
					}
				}
			}
		})
	})
}
