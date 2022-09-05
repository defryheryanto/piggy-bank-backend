package auth

import "github.com/defryheryanto/piggy-bank-backend/internal/errors"

var ErrInvalidCredential = errors.NewBadRequestError("Invalid Username or Password", "credentials entered not valid")
var ErrUsernameHasTaken = errors.NewBadRequestError("Username is already taken", "username is exists in database")
var ErrUserNotFound = errors.NewNotFoundError("User not found", "user not found")
