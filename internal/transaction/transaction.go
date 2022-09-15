package transaction

import (
	"context"
	"time"

	"github.com/defryheryanto/piggy-bank-backend/internal/storage"
	"golang.org/x/exp/slices"
)

type TransactionType string

const (
	IncomeType   TransactionType = "income"
	ExpenseType  TransactionType = "expense"
	SavingsType  TransactionType = "savings"
	TransferType TransactionType = "transfer"
)

var AvailableTransactionTypes = []TransactionType{IncomeType, ExpenseType, SavingsType, TransferType}

type Transaction struct {
	TransactionId   int       `gorm:"primaryKey;autoIncrement;column:transaction_id" json:"transaction_id"`
	UserId          int       `gorm:"column:user_id" json:"user_id"`
	AccountId       int       `gorm:"column:account_id" json:"account_id"`
	CategoryId      int       `gorm:"column:category_id" json:"category_id"`
	TransactionDate time.Time `gorm:"column:transaction_date" json:"transaction_date"`
	Description     string    `gorm:"column:description" json:"description"`
	Notes           string    `gorm:"column:notes" json:"notes"`
	Amount          float64   `gorm:"column:amount" json:"amount"`
	TransactionType string    `gorm:"column:transaction_type" json:"transaction_type"`
}

func (Transaction) TableName() string {
	return "transactions"
}

type CreateBasicTransactionPayload struct {
	UserId          int                         `json:"user_id"`
	AccountId       int                         `json:"account_id"`
	CategoryId      int                         `json:"category_id"`
	TransactionDate time.Time                   `json:"transaction_date"`
	Description     string                      `json:"description"`
	Notes           string                      `json:"notes"`
	Amount          float64                     `json:"amount"`
	TransactionType TransactionType             `json:"transaction_type"`
	Participants    []*CreateParticipantPayload `json:"participants"`
}

func (p *CreateBasicTransactionPayload) Validate() error {
	if p.UserId == 0 {
		return ErrInvalidUser
	}
	if p.AccountId == 0 {
		return ErrInvalidAccount
	}
	if p.CategoryId == 0 {
		return ErrInvalidCategory
	}
	if p.Amount < 0 {
		return ErrInvalidAmount
	}
	if !slices.Contains(AvailableTransactionTypes, p.TransactionType) {
		return ErrInvalidTransactionType
	}

	return nil
}

type TransactionRepository interface {
	Create(ctx context.Context, payload *Transaction) error
}

type TransactionService struct {
	repository         TransactionRepository
	participantService *ParticipantService
	manager            storage.Manager
}

func NewTransactionService(repo TransactionRepository, participantService *ParticipantService, manager storage.Manager) *TransactionService {
	return &TransactionService{repo, participantService, manager}
}

// Basic Transaction is transaction with type income or expense
func (s *TransactionService) CreateBasic(ctx context.Context, payload *CreateBasicTransactionPayload) error {
	err := payload.Validate()
	if err != nil {
		return err
	}
	trx := &Transaction{
		UserId:          payload.UserId,
		AccountId:       payload.AccountId,
		CategoryId:      payload.CategoryId,
		TransactionDate: payload.TransactionDate,
		Description:     payload.Description,
		Notes:           payload.Notes,
		Amount:          payload.Amount,
		TransactionType: string(payload.TransactionType),
	}

	err = s.manager.RunInTransaction(ctx, func(ctx context.Context) error {
		err = s.repository.Create(ctx, trx)
		if err != nil {
			return err
		}

		if len(payload.Participants) > 0 {
			err = s.participantService.BulkCreate(ctx, trx.TransactionId, payload.Participants)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
