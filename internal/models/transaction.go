package models

type Transaction struct {
	ID          int     `json:"id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	BudgetId    int     `json:"budget_id"`
	CategoryId  int     `json:"category_id"`
	Date        string  `json:"date"`
	Type        string  `json:"type"`
}
