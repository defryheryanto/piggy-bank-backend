package category

import (
	"github.com/defryheryanto/piggy-bank-backend/internal/errors"
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

type UpdateCategoryPayload struct {
	CategoryId   int          `json:"category_id"`
	CategoryName string       `json:"category_name"`
	CategoryType CategoryType `json:"category_type"`
	Budget       *int64       `json:"budget"`
	UserId       int64        `json:"user_id"`
}

func (p *UpdateCategoryPayload) Validate() error {
	if p.CategoryId == 0 {
		return errors.NewBadRequestError("Please fill Category ID", "category_id is required")
	}
	if p.CategoryName == "" {
		return errors.NewBadRequestError("Please fill Category Name", "category_name is required")
	}
	if p.CategoryType == "" {
		return errors.NewBadRequestError("Please fill Category Type", "category_type is required")
	}

	return nil
}

type CategoryRepository interface {
	Create(*Category) error
	GetByTypeAndUserId(categoryType CategoryType, userId int) []*Category
	GetById(categoryId int) *Category
	UpdateById(categoryId int, payload *Category) error
	DeleteById(categoryId int) error
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

func (s *CategoryService) UpdateCategory(payload *UpdateCategoryPayload) error {
	if payload.CategoryId == 0 {
		return ErrCategoryNotFound
	}

	existing := s.repository.GetById(payload.CategoryId)
	if existing == nil || existing.UserId != int(payload.UserId) {
		return ErrCategoryNotFound
	}

	if payload.CategoryName != "" {
		existing.CategoryName = payload.CategoryName
	}
	if payload.CategoryType != "" {
		err := ValidateCategoryType(payload.CategoryType)
		if err != nil {
			return err
		}
		existing.CategoryType = payload.CategoryType
	}
	if payload.Budget != nil {
		existing.Budget = *payload.Budget
	}

	err := s.repository.UpdateById(existing.CategoryId, existing)
	if err != nil {
		return err
	}

	return nil
}

func (s *CategoryService) DeleteById(categoryId, userId int) error {
	existing := s.repository.GetById(categoryId)
	if existing == nil || existing.UserId != userId {
		return ErrCategoryNotFound
	}

	err := s.repository.DeleteById(existing.CategoryId)
	if err != nil {
		return err
	}

	return nil
}

func (s *CategoryService) GetByCategoryTypeAndUserId(categoryType CategoryType, userId int) []*Category {
	return s.repository.GetByTypeAndUserId(categoryType, userId)
}

func ValidateCategoryType(categoryType CategoryType) error {
	isContain := slices.Contains(CategoryTypes, CategoryType(categoryType))
	if !isContain {
		return ErrInvalidCategoryType
	}

	return nil
}
