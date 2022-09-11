package transaction

import "github.com/defryheryanto/piggy-bank-backend/internal/errors"

var ErrEmptyParticipantName = errors.NewBadRequestError("Please fill Participant Name", "participant name is required")
var ErrInvalidAmount = errors.NewBadRequestError("Amount must be greater than 0", "amount is invalid")
