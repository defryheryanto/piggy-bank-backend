package storage

import (
	"strings"

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

func (s *AccountStorage) Update(id int, payload *account.Account) error {
	result := s.db.Where("account_id = ?", id).Updates(&payload)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *AccountStorage) GetByIdAndUser(accountId, userId int) *account.Account {
	var data *account.Account
	s.db.Where("account_id = ? AND user_id = ?", accountId, userId).First(&data)
	if data.AccountID == 0 {
		return nil
	}

	return data
}

func (s *AccountStorage) DeleteById(accountId int) error {
	result := s.db.Model(&account.Account{}).Delete("account_id = ?", accountId)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *AccountStorage) GetByFilter(filter *account.AccountFilter) []*account.Account {
	result := []*account.Account{}
	if filter == nil {
		s.db.Find(&result)
		return result
	}

	whereQueries := []string{}
	values := []interface{}{}
	if filter.UserID != 0 {
		whereQueries = append(whereQueries, "user_id = ?")
		values = append(values, filter.UserID)
	}
	if filter.AccountTypeId != 0 {
		whereQueries = append(whereQueries, "account_type_id = ?")
		values = append(values, filter.AccountTypeId)
	}
	if filter.ExcludedAccountTypeId != 0 {
		whereQueries = append(whereQueries, "account_type_id != ?")
		values = append(values, filter.ExcludedAccountTypeId)
	}

	s.db.Where(strings.Join(whereQueries, " AND "), values...).Find(&result)

	return result
}
