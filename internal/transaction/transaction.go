package transaction

import (
	"context"
	"time"

	"github.com/defryheryanto/piggy-bank-backend/internal/storage"
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

type TransferDetail struct {
	TransferDetailId int `gorm:"primaryKey;autoIncrement;column:transfer_detail_id" json:"transfer_detail_id"`
	TransactionId    int `gorm:"column:transaction_id" json:"transaction_id"`
	TargetAccountId  int `gorm:"column:target_account_id" json:"target_account_id"`
}

func (TransferDetail) TableName() string {
	return "transfer_details"
}

type SavingDetail struct {
	SavingDetailId  int `gorm:"primaryKey;autoIncrement;column:saving_detail_id" json:"saving_detail_id"`
	TransactionId   int `gorm:"column:transaction_id" json:"transaction_id"`
	TargetAccountId int `gorm:"column:target_account_id" json:"target_account_id"`
}

func (SavingDetail) TableName() string {
	return "saving_details"
}

type TransactionRepository interface {
	Create(ctx context.Context, payload *Transaction) error
	CreateTransferDetail(ctx context.Context, payload *TransferDetail) error
	CreateSavingDetail(ctx context.Context, payload *SavingDetail) error
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

func (s *TransactionService) CreateTransfer(ctx context.Context, payload *CreateTransferPayload) error {
	err := payload.Validate()
	if err != nil {
		return err
	}

	trx := &Transaction{
		UserId:          payload.UserId,
		AccountId:       payload.SourceAccountId,
		TransactionDate: payload.TransactionDate,
		Description:     payload.Description,
		Notes:           payload.Notes,
		Amount:          payload.Amount,
		TransactionType: string(TransferType),
	}

	err = s.manager.RunInTransaction(ctx, func(ctx context.Context) error {
		err = s.repository.Create(ctx, trx)
		if err != nil {
			return err
		}

		transferDetail := &TransferDetail{
			TransactionId:   trx.TransactionId,
			TargetAccountId: payload.TargetAccountId,
		}

		err = s.repository.CreateTransferDetail(ctx, transferDetail)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *TransactionService) CreateSaving(ctx context.Context, payload *CreateSavingPayload) error {
	err := payload.Validate()
	if err != nil {
		return err
	}

	trx := &Transaction{
		UserId:          payload.UserId,
		AccountId:       payload.SourceAccountId,
		TransactionDate: payload.TransactionDate,
		Description:     payload.Description,
		Notes:           payload.Notes,
		Amount:          payload.Amount,
		TransactionType: string(SavingsType),
	}

	err = s.manager.RunInTransaction(ctx, func(ctx context.Context) error {
		err = s.repository.Create(ctx, trx)
		if err != nil {
			return err
		}

		savingDetail := &SavingDetail{
			TransactionId:   trx.TransactionId,
			TargetAccountId: payload.TargetAccountId,
		}

		err = s.repository.CreateSavingDetail(ctx, savingDetail)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
