package storage

import (
	"github.com/defryheryanto/piggy-bank-backend/internal/account"
	"gorm.io/gorm"
)

type AccountStorage struct {
	db *gorm.DB
}

func NewAccountStorage(db *gorm.DB) *AccountStorage {
	return &AccountStorage{db}
}

func (s *AccountStorage) GetTypes() []*account.AccountType {
	var datas []*account.AccountType
	s.db.Find(&datas)
	return datas
}

func (s *AccountStorage) GetTypeById(typeId int) *account.AccountType {
	var data *account.AccountType
	s.db.Where("account_type_id = ?", typeId).First(&data)
	if data.AccountTypeID == 0 {
		return nil
	}
	return data
}

func (s *AccountStorage) GetByUserIdAndType(userID, typeID int) []*account.Account {
	var datas []*account.Account
	s.db.Where("user_id = ? AND account_type_id = ?", userID, typeID).Find(&datas)
	return datas
}

func (s *AccountStorage) Create(payload *account.Account) error {
	result := s.db.Create(payload)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
