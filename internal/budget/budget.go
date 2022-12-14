package budget

import (
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

type ActiveBudgetDetail struct {
	CategoryId   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
	Budget       int64  `json:"budget"`
}

type BudgetRepository interface {
	Create(payload *Budget) error
	Update(payload *Budget) error
	GetByMonthAndYear(categoryId, month, year int) *Budget
}

type BudgetService struct {
	repository      BudgetRepository
	categoryService *category.CategoryService
}

func NewBudgetService(repo BudgetRepository, category *category.CategoryService) *BudgetService {
	return &BudgetService{repo, category}
}

func (s *BudgetService) CreateOrUpdate(payload *CreateBudgetPayload) error {
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

	existing := s.repository.GetByMonthAndYear(budget.CategoryId, budget.Month, budget.Year)
	if existing != nil {
		budget.BudgetId = existing.BudgetId
		err = s.repository.Update(budget)
	} else {
		err = s.repository.Create(budget)
	}

	if err != nil {
		return err
	}

	return nil
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

func (s *BudgetService) GetActiveBudget(userId, month, year int) []*ActiveBudgetDetail {
	categories := s.categoryService.GetByCategoryTypeAndUserId(category.ExpenseCategory, userId)
	activeBudgets := []*ActiveBudgetDetail{}
	for _, cat := range categories {
		bdg := s.repository.GetByMonthAndYear(cat.CategoryId, month, year)
		active := cat.Budget
		if bdg != nil {
			active = bdg.Budget
		}
		activeBudgets = append(activeBudgets, &ActiveBudgetDetail{
			CategoryId:   cat.CategoryId,
			CategoryName: cat.CategoryName,
			Budget:       active,
		})
	}

	return activeBudgets
}
