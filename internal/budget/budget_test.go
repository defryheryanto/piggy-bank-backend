package budget_test

import (
	"testing"

	"github.com/defryheryanto/piggy-bank-backend/internal/budget"
	budget_storage "github.com/defryheryanto/piggy-bank-backend/internal/budget/sql"
	"github.com/defryheryanto/piggy-bank-backend/test"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupService(db *gorm.DB) *budget.BudgetService {
	budgetStorage := budget_storage.NewBudgetStorage(db)

	return budget.NewBudgetService(budgetStorage)
}

func setupDatabase(t *testing.T) *gorm.DB {
	db := test.SetupTestDatabase(t, "../../.env", "../../db/migrations")
	return db
}

func TestCreate(t *testing.T) {
	db := setupDatabase(t)

	test.RunTestWithDB(db, t, func(t *testing.T, db *gorm.DB) {
		service := setupService(db)

		t.Run("should insert to db", func(t *testing.T) {
			payload := &budget.CreateBudgetPayload{
				CategoryId: 1,
				Month:      9,
				Year:       2022,
				Budget:     1500000,
			}

			err := service.Create(payload)
			assert.NoError(t, err)
		})

		t.Run("return error if month is invalid", func(t *testing.T) {
			payload := &budget.CreateBudgetPayload{
				CategoryId: 2,
				Month:      0,
				Year:       2022,
				Budget:     1500000,
			}

			err := service.Create(payload)
			assert.ErrorIs(t, err, budget.ErrInvalidMonthBudget)

			payload = &budget.CreateBudgetPayload{
				CategoryId: 2,
				Month:      13,
				Year:       2022,
				Budget:     1500000,
			}
			err = service.Create(payload)
			assert.ErrorIs(t, err, budget.ErrInvalidMonthBudget)
		})
	})
}
