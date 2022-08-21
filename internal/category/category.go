package category

import (
	"golang.org/x/exp/slices"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type CategoryType string

func (ct *CategoryType) DisplayName() string {
	caser := cases.Title(language.English)
	return caser.String(string(*ct))
}

const (
	IncomeCategory  CategoryType = "income"
	ExpenseCategory CategoryType = "expense"
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

type CategoryTypeDetail struct {
	CategoryType        CategoryType `json:"category_type"`
	CategoryTypeDisplay string       `json:"category_type_display"`
	Categories          []*Category  `json:"categories"`
}

type CategoryRepository interface {
	Create(*Category) error
	GetByTypeAndUserId(categoryType CategoryType, userId int) []*Category
	GetById(categoryId int) *Category
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

func (s *CategoryService) GetCategoryTypeDetails(userId int) []*CategoryTypeDetail {
	result := []*CategoryTypeDetail{}
	for _, categoryType := range CategoryTypes {
		categories := s.repository.GetByTypeAndUserId(categoryType, userId)
		result = append(result, &CategoryTypeDetail{
			CategoryType:        categoryType,
			CategoryTypeDisplay: categoryType.DisplayName(),
			Categories:          categories,
		})
	}

	return result
}

func (s *CategoryService) GetCategoryById(categoryId int) (*Category, error) {
	cat := s.repository.GetById(categoryId)
	if cat == nil {
		return nil, ErrCategoryNotFound
	}

	return cat, nil
}

func ValidateCategoryType(categoryType CategoryType) error {
	isContain := slices.Contains(CategoryTypes, CategoryType(categoryType))
	if !isContain {
		return ErrInvalidCategoryType
	}

	return nil
}
