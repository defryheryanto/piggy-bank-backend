package budget

type CreateBudgetPayload struct {
	CategoryId int   `json:"category_id"`
	Month      int   `json:"month"`
	Year       int   `json:"year"`
	Budget     int64 `json:"budget"`
}

func (p *CreateBudgetPayload) Validate() error {
	if p.Month < 1 || p.Month > 12 {
		return ErrInvalidMonthBudget
	}

	return nil
}
