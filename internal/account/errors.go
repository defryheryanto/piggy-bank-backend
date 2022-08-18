package account

import "github.com/defryheryanto/piggy-bank-backend/internal/errors"

var ErrAccountTypeNotFound = errors.NewNotFoundError("Account Type not found", "account type does not exists")
var ErrAccountNotFound = errors.NewNotFoundError("Account not found", "account does not exists")
