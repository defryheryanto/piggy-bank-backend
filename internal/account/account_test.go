package account_test

import (
	"testing"

	"github.com/defryheryanto/piggy-bank-backend/internal/account"
	"github.com/defryheryanto/piggy-bank-backend/internal/account/storage"
	"github.com/defryheryanto/piggy-bank-backend/test"
	"github.com/stretchr/testify/assert"
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

func TestCreateAccount(t *testing.T) {
	db := test.SetupTestDatabase(t, "../../.env", "../../db/migrations")
	test.RunTestWithDB(db, t, func(t *testing.T, db *gorm.DB) {
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

		t.Run("should insert to database", func(t *testing.T) {
			service := setupAccountService(db)
			payload := &account.Account{
				AccountName:   "BCA",
				AccountTypeID: 1,
				UserID:        1,
			}

			err := service.CreateAccount(payload)
			if err != nil {
				t.Errorf("should have no error when creating account")
			}

			var insertedData *account.Account
			db.Where("account_name = ?", payload.AccountName).First(&insertedData)
			assert.Equal(t, payload.AccountTypeID, insertedData.AccountTypeID)
			assert.Equal(t, payload.UserID, payload.UserID)

			if insertedData.CreatedAt.IsZero() {
				t.Errorf("created_at should not be nil")
			}
			if insertedData.UpdatedAt.IsZero() {
				t.Errorf("updated_at should not be nil")
			}
		})
	})
}
