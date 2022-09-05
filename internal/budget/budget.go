package budget

import (
	"fmt"
	"time"

	"github.com/defryheryanto/piggy-bank-backend/internal/category"
)

type Budget struct {
	BudgetId   int   `gorm:"primaryKey;autoIncrement;column:budget_id" json:"budget_id"`
	CategoryId int   `gorm:"column:category_id" json:"category_id"`
	Month      int   `gorm:"column:month" json:"month"`
	Year       int   `gorm:"column:year" json:"year"`
	Budget     int64 `gorm:"column:budget" json:"budget"`
}

func (Budget) TableName() string {
	return "budgets"
}

type BudgetDetail struct {
	Month  string `json:"month"`
	Year   int    `json:"year"`
	Budget int64  `json:"budget"`
}

type BudgetSummary struct {
	DefaultBudget int64           `json:"default_budget"`
	Budgets       []*BudgetDetail `json:"budgets"`
}

type BudgetRepository interface {
	Create(payload *Budget) error
	GetByMonthAndYear(categoryId, month, year int) *Budget
}

type BudgetService struct {
	repository      BudgetRepository
	categoryService *category.CategoryService
}

func NewBudgetService(repo BudgetRepository, category *category.CategoryService) *BudgetService {
	return &BudgetService{repo, category}
}

func (s *BudgetService) Create(payload *CreateBudgetPayload) error {
	err := payload.Validate()
	if err != nil {
		return err
	}

	budget := &Budget{
		CategoryId: payload.CategoryId,
		Month:      payload.Month,
		Year:       payload.Year,
		Budget:     payload.Budget,
	}

	return s.repository.Create(budget)
}

func (s *BudgetService) GetBudgetYearSummary(categoryId, year int) (*BudgetSummary, error) {
	category, err := s.categoryService.GetCategoryById(categoryId)
	if err != nil {
		return nil, err
	}

	budgets := []*BudgetDetail{}
	for i := 1; i <= 12; i++ {
		budget := s.repository.GetByMonthAndYear(categoryId, i, year)
		value := category.Budget
		if budget != nil {
			fmt.Println(*budget)
			value = budget.Budget
		}

		budgets = append(budgets, &BudgetDetail{
			Month:  time.Month(i).String(),
			Year:   year,
			Budget: value,
		})
	}

	return &BudgetSummary{
		DefaultBudget: category.Budget,
		Budgets:       budgets,
	}, nil
}
