package handler

import "github.com/defryheryanto/piggy-bank-backend/internal/errors"

type CreateCategoryPayload struct {
	CategoryName string `json:"name"`
	CategoryType string `json:"type"`
	Budget       int64  `json:"budget"`
}

func (p *CreateCategoryPayload) Validate() error {
	if p.CategoryName == "" {
		return errors.NewBadRequestError("Please fill name", "category name is required")
	}
	if p.CategoryType == "" {
		return errors.NewBadRequestError("Please fill type", "category type is required")
	}

	return nil
}
