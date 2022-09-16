package account_test

import (
	"testing"

	"github.com/defryheryanto/piggy-bank-backend/internal/account"
	storage "github.com/defryheryanto/piggy-bank-backend/internal/account/sql"
	"github.com/defryheryanto/piggy-bank-backend/test"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupAccountService(db *gorm.DB) *account.AccountService {
	accountStorage := storage.NewAccountStorage(db)
	return account.NewAccountService(accountStorage)
}

func populateAccountType(db *gorm.DB) {
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

		t.Run("return error if account type is not exists", func(t *testing.T) {
			service := setupAccountService(db)
			payload := &account.Account{
				AccountName:   "BCA",
				AccountTypeID: 99,
				UserID:        1,
			}

			err := service.CreateAccount(payload)
			assert.NotNil(t, err)
			assert.ErrorIs(t, err, account.ErrAccountTypeNotFound)
		})
	})
}

func TestUpdateAccount(t *testing.T) {
	db := test.SetupTestDatabase(t, "../../.env", "../../db/migrations")
	test.RunTestWithDB(db, t, func(t *testing.T, db *gorm.DB) {
		t.Run("should update data", func(t *testing.T) {
			populateAccountType(db)

			acc := &account.Account{
				AccountID:     1,
				AccountName:   "BCA",
				AccountTypeID: 1,
				UserID:        1,
			}

			db.Create(acc)

			service := setupAccountService(db)

			payload := &account.UpdateAccountPayload{
				AccountID:     1,
				AccountName:   "Mandiri",
				AccountTypeID: 2,
				UserID:        1,
			}

			err := service.UpdateAccount(payload)
			assert.NoError(t, err)

			var inserted *account.Account
			db.Where("account_id = ?", acc.AccountID).First(&inserted)

			if inserted.AccountName != payload.AccountName {
				t.Errorf("account_name not updated")
			}
			if inserted.AccountTypeID != payload.AccountTypeID {
				t.Errorf("account_type_id not updated")
			}
		})

		t.Run("return error if account not found", func(t *testing.T) {
			service := setupAccountService(db)

			payload := &account.UpdateAccountPayload{
				AccountID:     99,
				AccountName:   "Mandiri",
				AccountTypeID: 2,
				UserID:        1,
			}

			err := service.UpdateAccount(payload)
			assert.ErrorIs(t, err, account.ErrAccountNotFound)
		})

		t.Run("return error if payload account id is 0 or nil", func(t *testing.T) {
			service := setupAccountService(db)

			payload := &account.UpdateAccountPayload{
				AccountID:     0,
				AccountName:   "Mandiri",
				AccountTypeID: 99,
				UserID:        1,
			}

			err := service.UpdateAccount(payload)
			assert.ErrorIs(t, err, account.ErrAccountNotFound)

			err = service.UpdateAccount(nil)
			assert.ErrorIs(t, err, account.ErrAccountNotFound)
		})

		t.Run("return error if account user id is different with payload", func(t *testing.T) {
			populateAccountType(db)

			acc := &account.Account{
				AccountID:     1,
				AccountName:   "BCA",
				AccountTypeID: 1,
				UserID:        1,
			}

			db.Create(acc)

			service := setupAccountService(db)

			payload := &account.UpdateAccountPayload{
				AccountID:     1,
				AccountName:   "Mandiri",
				AccountTypeID: 2,
				UserID:        2,
			}

			err := service.UpdateAccount(payload)
			assert.ErrorIs(t, err, account.ErrAccountNotFound)
		})
	})
}

func TestDeleteAccount(t *testing.T) {
	db := test.SetupTestDatabase(t, "../../.env", "../../db/migrations")
	test.RunTestWithDB(db, t, func(t *testing.T, db *gorm.DB) {
		service := setupAccountService(db)
		t.Run("should delete the data", func(t *testing.T) {
			acc := &account.Account{
				AccountID:     1,
				AccountName:   "BCA",
				AccountTypeID: 1,
				UserID:        1,
			}
			db.Create(&acc)

			err := service.DeleteById(acc.AccountID, acc.UserID)
			assert.NoError(t, err)
		})

		t.Run("return error if account not found", func(t *testing.T) {
			err := service.DeleteById(99, 1)
			assert.ErrorIs(t, err, account.ErrAccountNotFound)
		})

		t.Run("return error if user id is not match", func(t *testing.T) {
			acc := &account.Account{
				AccountID:     1,
				AccountName:   "BCA",
				AccountTypeID: 1,
				UserID:        1,
			}
			db.Create(&acc)

			err := service.DeleteById(acc.AccountID, 2)
			assert.ErrorIs(t, err, account.ErrAccountNotFound)
		})
	})
}

func TestGetList_FilterByUserId(t *testing.T) {
	db := test.SetupTestDatabase(t, "../../.env", "../../db/migrations")
	service := setupAccountService(db)

	tables := []string{
		account.Account{}.TableName(),
	}

	test.TruncateAfterTest(t, db, tables, func() {
		init := []*account.Account{
			{
				AccountName:   "BCA",
				AccountTypeID: 1,
				UserID:        1,
			},
			{
				AccountName:   "BCA",
				AccountTypeID: 2,
				UserID:        2,
			},
		}

		db.Create(&init)

		results := service.GetList(&account.AccountFilter{
			UserID: 1,
		})
		assert.Equal(t, 1, len(results))
	})
}

func TestGetList_FilterByAccountTypeId(t *testing.T) {
	db := test.SetupTestDatabase(t, "../../.env", "../../db/migrations")
	service := setupAccountService(db)

	tables := []string{
		account.Account{}.TableName(),
	}

	test.TruncateAfterTest(t, db, tables, func() {
		init := []*account.Account{
			{
				AccountName:   "BCA",
				AccountTypeID: 1,
				UserID:        1,
			},
			{
				AccountName:   "Blu",
				AccountTypeID: 2,
				UserID:        1,
			},
		}

		db.Create(&init)

		results := service.GetList(&account.AccountFilter{
			AccountTypeId: 1,
		})
		assert.Equal(t, 1, len(results))
	})
}

func TestGetList_FilterByExcludedAccountTypeId(t *testing.T) {
	db := test.SetupTestDatabase(t, "../../.env", "../../db/migrations")
	service := setupAccountService(db)

	tables := []string{
		account.Account{}.TableName(),
	}

	test.TruncateAfterTest(t, db, tables, func() {
		init := []*account.Account{
			{
				AccountName:   "BCA",
				AccountTypeID: 1,
				UserID:        1,
			},
			{
				AccountName:   "Blu",
				AccountTypeID: 2,
				UserID:        1,
			},
			{
				AccountName:   "Blu",
				AccountTypeID: 1,
				UserID:        1,
			},
		}

		db.Create(&init)

		results := service.GetList(&account.AccountFilter{
			ExcludedAccountTypeId: 2,
		})
		assert.Equal(t, 2, len(results))
	})
}
