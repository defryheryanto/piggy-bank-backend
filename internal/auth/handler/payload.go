package handler

import "github.com/defryheryanto/piggy-bank-backend/internal/errors"

type registerPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (p *registerPayload) Validate() error {
	if p.Username == "" {
		return errors.NewBadRequestError("Please fill username", "username is empty")
	}
	if p.Password == "" {
		return errors.NewBadRequestError("Please fill password", "password is empty")
	}

	return nil
}

type loginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (p *loginPayload) Validate() error {
	if p.Username == "" {
		return errors.NewBadRequestError("Please fill username", "username is empty")
	}
	if p.Password == "" {
		return errors.NewBadRequestError("Please fill password", "password is empty")
	}

	return nil
}
