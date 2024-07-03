package models

type Budget struct {
	ID          int     `json:"id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	CategoryId  int     `json:"category_id"`
	UserId      int     `json:"user_id"`
}
