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

	rows, err := h.DB.Query("SELECT id, name, amount FROM budgets")
	if err != nil {
		log.Printf("Error fetching budgets: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var budget models.Budget
		if err := rows.Scan(&budget.ID, &budget.Description, &budget.Amount); err != nil {
			log.Printf("Error scanning budget row: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		budgets = append(budgets, budget)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(budgets)
}

func (h *BudgetHandler) CreateBudget(w http.ResponseWriter, r *http.Request) {
	var budget models.Budget

	if err := json.NewDecoder(r.Body).Decode(&budget); err != nil {
		log.Printf("Error decoding budget: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Insert the new budget into the database
	result, err := h.DB.Exec("INSERT INTO budgets (description, amount, category_id, user_id) VALUES (?, ?, ?, 1)", budget.Description, budget.Amount, budget.CategoryId)
	if err != nil {
		log.Printf("Error inserting budget: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Retrieve the ID of the newly inserted budget
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the ID of the budget and return the created budget
	budget.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(budget)
}
