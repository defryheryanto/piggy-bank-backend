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
	UserID        int       `gorm:"column:user_id"`
	Balance       int64     `gorm:"column:balance"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Account) TableName() string {
	return "accounts"
}

type AccountSummary struct {
	*AccountType
	Accounts []*Account `json:"accounts"`
}

type UpdateAccountPayload struct {
	AccountID     int    `json:"account_id"`
	AccountName   string `json:"account_name"`
	AccountTypeID int    `json:"account_type_id"`
	UserID        int    `json:"user_id"`
}

type AccountRepository interface {
	GetTypes() []*AccountType
	GetTypeById(int) *AccountType
	GetByUserIdAndType(userID, typeID int) []*Account
	Create(*Account) error
	Update(id int, payload *Account) error
	GetByIdAndUser(accountId, userId int) *Account
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

func (s *AccountService) GetAccountsByUser(userID int) []*AccountSummary {
	types := s.accountStorage.GetTypes()
	result := []*AccountSummary{}
	for _, t := range types {
		accounts := s.accountStorage.GetByUserIdAndType(userID, t.AccountTypeID)
		result = append(result, &AccountSummary{
			AccountType: t,
			Accounts:    accounts,
		})
	}

	return result
}

func (s *AccountService) CreateAccount(payload *Account) error {
	accType := s.accountStorage.GetTypeById(payload.AccountTypeID)
	if accType == nil {
		return ErrAccountTypeNotFound
	}

	err := s.accountStorage.Create(payload)
	if err != nil {
		return err
	}

	return nil
}

func (s *AccountService) UpdateAccount(payload *UpdateAccountPayload) error {
	if payload == nil || payload.AccountID == 0 {
		return ErrAccountTypeNotFound
	}

	currentAccount := s.accountStorage.GetByIdAndUser(payload.UserID, payload.AccountID)
	if currentAccount == nil {
		return ErrAccountTypeNotFound
	}

	if payload.AccountName != "" {
		currentAccount.AccountName = payload.AccountName
	}
	if payload.AccountTypeID != 0 {
		currentAccount.AccountTypeID = payload.AccountTypeID
	}
	currentAccount.UpdatedAt = time.Now()

	err := s.accountStorage.Update(payload.AccountID, currentAccount)
	if err != nil {
		return err
	}

	return nil
}
