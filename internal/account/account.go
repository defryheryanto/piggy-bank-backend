package account

import "time"

type AccountType struct {
	AccountTypeID   int    `gorm:"primaryKey;autoIncrement;column:account_type_id" json:"account_type_id"`
	AccountTypeName string `gorm:"column:account_type_name" json:"account_type_name"`
}

func (AccountType) TableName() string {
	return "account_types"
}

type Account struct {
	AccountID     int       `gorm:"primaryKey;autoIncrement;column:account_id" json:"account_id"`
	AccountName   string    `gorm:"column:account_name" json:"account_name"`
	AccountTypeID int       `gorm:"column:account_type_id" json:"account_type_id"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Account) TableName() string {
	return "accounts"
}

type AccountRepository interface {
	GetTypes() []*AccountType
}

type AccountService struct {
	accountStorage AccountRepository
}

func NewAccountService(accountRepository AccountRepository) *AccountService {
	return &AccountService{accountRepository}
}

func (s *AccountService) GetTypes() []*AccountType {
	return s.accountStorage.GetTypes()
}
