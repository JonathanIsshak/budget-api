// internal/handlers/budget_handler.go

package handlers

import (
	"budgeting-app/internal/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type BudgetHandler struct {
	DB *sql.DB
}

func NewBudgetHandler(db *sql.DB) *BudgetHandler {
	return &BudgetHandler{DB: db}
}

func (h *BudgetHandler) GetBudgets(w http.ResponseWriter, r *http.Request) {
	var budgets []models.Budget

	// Example query
	rows, err := h.DB.Query("SELECT id, name, amount FROM budgets")
	if err != nil {
		log.Printf("Error fetching budgets: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var budget models.Budget
		if err := rows.Scan(&budget.ID, &budget.Name, &budget.Amount); err != nil {
			log.Printf("Error scanning budget row: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		budgets = append(budgets, budget)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(budgets)
}

// Add more handlers for CRUD operations on budgets
