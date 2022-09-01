package budget

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

type BudgetRepository interface {
	Create(payload *Budget) error
}

type BudgetService struct {
	repository BudgetRepository
}

func NewBudgetService(repo BudgetRepository) *BudgetService {
	return &BudgetService{repo}
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
