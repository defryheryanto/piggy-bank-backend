package transaction

import "github.com/defryheryanto/piggy-bank-backend/internal/errors"

var ErrEmptyParticipantName = errors.NewBadRequestError("Please fill Participant Name", "participant name is required")
var ErrInvalidAmount = errors.NewBadRequestError("Amount must be greater than 0", "amount is invalid")
var ErrInvalidUser = errors.NewBadRequestError("User is invalid to create transaction", "user is invalid")
var ErrInvalidAccount = errors.NewBadRequestError("Account is invalid", "account is invalid")
var ErrInvalidCategory = errors.NewBadRequestError("Category is invalid", "category is invalid")
var ErrInvalidTransactionType = errors.NewBadRequestError("Transaction Type is invalid", "transaction type is invalid")
var ErrInvalidSourceAccount = errors.NewBadRequestError("Source Account is invalid", "source account is invalid")
var ErrInvalidTargetAccount = errors.NewBadRequestError("Target Account is invalid", "target account is invalid")
var ErrInvalidTransactionDate = errors.NewBadRequestError("Transaction Date is invalid", "transaction date is invalid")
