package transaction

import (
	"time"

	"golang.org/x/exp/slices"
)

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

type CreateTransferPayload struct {
	UserId          int       `json:"user_id"`
	SourceAccountId int       `json:"source_account_id"`
	TargetAccountId int       `json:"target_account_id"`
	TransactionDate time.Time `json:"transaction_date"`
	Description     string    `json:"description"`
	Notes           string    `json:"notes"`
	Amount          float64   `json:"amount"`
}

func (p *CreateTransferPayload) Validate() error {
	if p.UserId == 0 {
		return ErrInvalidUser
	}
	if p.SourceAccountId == 0 {
		return ErrInvalidSourceAccount
	}
	if p.TargetAccountId == 0 {
		return ErrInvalidTargetAccount
	}
	if p.TransactionDate.IsZero() {
		return ErrInvalidTransactionDate
	}
	if p.Amount < 0 {
		return ErrInvalidAmount
	}

	return nil
}

type CreateSavingPayload struct {
	UserId          int       `json:"user_id"`
	SourceAccountId int       `json:"source_account_id"`
	TargetAccountId int       `json:"target_account_id"`
	TransactionDate time.Time `json:"transaction_date"`
	Description     string    `json:"description"`
	Notes           string    `json:"notes"`
	Amount          float64   `json:"amount"`
}

func (p *CreateSavingPayload) Validate() error {
	if p.UserId == 0 {
		return ErrInvalidUser
	}
	if p.SourceAccountId == 0 {
		return ErrInvalidSourceAccount
	}
	if p.TargetAccountId == 0 {
		return ErrInvalidTargetAccount
	}
	if p.TransactionDate.IsZero() {
		return ErrInvalidTransactionDate
	}
	if p.Amount < 0 {
		return ErrInvalidAmount
	}

	return nil
}
