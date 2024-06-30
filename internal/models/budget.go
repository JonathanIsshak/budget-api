// internal/models/budget.go

package models

type Budget struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

// Add more models as needed
