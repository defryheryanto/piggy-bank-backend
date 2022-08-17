package auth

import "github.com/defryheryanto/piggy-bank-backend/internal/errors"

var InvalidCredentialError = errors.NewBadRequestError("Invalid Username or Password", "credentials entered not valid")
var UsernameHasTakenError = errors.NewBadRequestError("Username is already taken", "username is exists in database")
