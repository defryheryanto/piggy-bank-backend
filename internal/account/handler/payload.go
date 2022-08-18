package handler

import "github.com/defryheryanto/piggy-bank-backend/internal/errors"

type CreateAccountPayload struct {
	AccountTypeId int    `json:"account_type_id"`
	Name          string `json:"name"`
}

func (p *CreateAccountPayload) Validate() error {
	if p.AccountTypeId == 0 {
		return errors.NewBadRequestError("Please fill account type id", "account_type_id is required")
	}
	if p.Name == "" {
		return errors.NewBadRequestError("Please fill name", "name is required")
	}

	return nil
}
