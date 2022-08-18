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
	test.RunTestWithDB(db, t, func(t *testing.T, db *gorm.DB) {
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
		db.Create(&types)

		t.Run("return all account types", func(t *testing.T) {
			existingTypes := service.GetTypes()
			if len(existingTypes) != len(types) {
				t.Errorf("type is not complete")
			}
		})
	})
}

func TestGetAccountsByUserAndType(t *testing.T) {
	db := test.SetupTestDatabase(t, "../../.env", "../../db/migrations")
	test.RunTestWithDB(db, t, func(t *testing.T, db *gorm.DB) {
		service := setupAccountService(db)

		types := []*account.AccountType{
			{
				AccountTypeID:   1,
				AccountTypeName: "Account",
			},
			{
				AccountTypeID:   2,
				AccountTypeName: "Credit Card",
			},
			{
				AccountTypeID:   3,
				AccountTypeName: "Savings",
			},
		}
		db.Create(&types)

		accounts := []*account.Account{
			{
				AccountName:   "BCA",
				AccountTypeID: types[0].AccountTypeID,
				UserID:        1,
			},
			{
				AccountName:   "Bibit",
				AccountTypeID: types[2].AccountTypeID,
				UserID:        1,
			},
		}
		db.Create(&accounts)

		t.Run("return accounts grouped by account type", func(t *testing.T) {
			accounts := service.GetAccountsByUser(1)
			if len(accounts) != len(types) {
				t.Errorf("accounts length should match account types count")
			}
		})
	})
}
