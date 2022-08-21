package category

import "github.com/defryheryanto/piggy-bank-backend/internal/errors"

var ErrInvalidCategoryType = errors.NewBadRequestError("Category Type is not valid", "invalid category type")
var ErrCategoryNotFound = errors.NewNotFoundError("Category not found", "category not found")
