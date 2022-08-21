package category

import "golang.org/x/exp/slices"

type CategoryType string

const (
	IncomeCategory  CategoryType = "income"
	ExpenseCategory CategoryType = "Expense"
)

var CategoryTypes = []CategoryType{
	IncomeCategory,
	ExpenseCategory,
}

type Category struct {
	CategoryId   int          `gorm:"primaryKey;autoIncremen;column:category_id" json:"category_id"`
	CategoryName string       `gorm:"column:category_name" json:"category_name"`
	CategoryType CategoryType `gorm:"column:category_type" json:"category_type"`
	UserId       int          `gorm:"column:user_id" json:"user_id"`
	Budget       int64        `gorm:"column:budget" json:"budget"`
}

func (Category) TableName() string {
	return "categories"
}

type CategoryRepository interface {
	Create(*Category) error
}

type CategoryService struct {
	repository CategoryRepository
}

func NewCategoryService(repo CategoryRepository) *CategoryService {
	return &CategoryService{repo}
}

func (s *CategoryService) Create(payload *Category) error {
	err := ValidateCategoryType(payload.CategoryType)
	if err != nil {
		return err
	}

	err = s.repository.Create(payload)
	if err != nil {
		return err
	}

	return nil
}

func ValidateCategoryType(categoryType CategoryType) error {
	isContain := slices.Contains(CategoryTypes, CategoryType(categoryType))
	if !isContain {
		return ErrInvalidCategoryType
	}

	return nil
}
