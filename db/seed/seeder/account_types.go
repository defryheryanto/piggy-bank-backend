package seeder

import (
	"log"

	"github.com/defryheryanto/piggy-bank-backend/internal/account"
	"gorm.io/gorm"
)

type AccountTypeSeeder struct {
	db *gorm.DB
}

func NewAccountTypeSeeder(db *gorm.DB) *AccountTypeSeeder {
	return &AccountTypeSeeder{db}
}

func (s *AccountTypeSeeder) Seed() {
	accountTypes := []*account.AccountType{
		{
			AccountTypeID:   1,
			AccountTypeName: "Account",
		},
		{
			AccountTypeID:   2,
			AccountTypeName: "Savings",
		},
	}

	for _, accountType := range accountTypes {
		var existing *account.AccountType
		s.db.Where("account_type_name = ?", accountType.AccountTypeName).Find(&existing)
		if existing.AccountTypeID == 0 {
			err := s.db.Create(&accountType)
			if err.Error != nil {
				log.Fatalf("error insert account type - %s: %v", accountType.AccountTypeName, err.Error)
			}
		}
	}
}
