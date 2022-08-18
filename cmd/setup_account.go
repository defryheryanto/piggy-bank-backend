package main

import (
	"github.com/defryheryanto/piggy-bank-backend/internal/account"
	storage "github.com/defryheryanto/piggy-bank-backend/internal/account/sql"
	"gorm.io/gorm"
)

func setupAccountService(db *gorm.DB) *account.AccountService {
	accountStorage := storage.NewAccountStorage(db)
	return account.NewAccountService(accountStorage)
}
