package budget

import "github.com/defryheryanto/piggy-bank-backend/internal/errors"

var ErrInvalidMonthBudget = errors.NewBadRequestError("Month should be numeric with value 1 to 12", "month numeric is invalid")
