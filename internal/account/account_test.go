package account_test

import (
	"testing"

	"github.com/defryheryanto/piggy-bank-backend/internal/account"
	"github.com/defryheryanto/piggy-bank-backend/internal/account/storage"
	"github.com/defryheryanto/piggy-bank-backend/test"
	"gorm.io/gorm"
)

func setupAccountService(db *gorm.DB) *account.AccountService {
	accountStorage := storage.NewAccountStorage(db)
	return account.NewAccountService(accountStorage)
}

func TestGetTypes(t *testing.T) {
	db := test.SetupTestDatabase(t, "../../.env", "../../db/migrations")
	db = db.Begin()
	service := setupAccountService(db)

	types := []*account.AccountType{
		{
			AccountTypeName: "Account",
		},
		{
			AccountTypeName: "Credit Card",
		},
		{
			AccountTypeName: "Savings",
		},
	}

	t.Run("return all account types", func(t *testing.T) {
		db.Create(&types)
		existingTypes := service.GetTypes()
		if len(existingTypes) != len(types) {
			t.Errorf("type is not complete")
		}
	})

	db.Rollback()
}
