package auth

import "github.com/defryheryanto/piggy-bank-backend/internal/errors"

var InvalidCredentialError = errors.NewBadRequestError("Invalid Username or Password", "credentials entered not valid")
