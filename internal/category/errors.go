package category

import "github.com/defryheryanto/piggy-bank-backend/internal/errors"

var ErrInvalidCategoryType = errors.NewBadRequestError("Category Type is not valid", "invalid category type")
