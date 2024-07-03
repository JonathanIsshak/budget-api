package handlers

import (
	"budgeting-app/internal/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type TransactionHandler struct {
	DB *sql.DB
}

func NewTransactionHandler(db *sql.DB) *TransactionHandler {
	return &TransactionHandler{DB: db}
}

func (h *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	var transactions []models.Transaction

	rows, err := h.DB.Query("SELECT id, description, amount, date, budget_id, category_id, type FROM transactions")
	if err != nil {
		log.Printf("Error fetching transactions: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.Description, &transaction.Amount, &transaction.Date, &transaction.BudgetId, &transaction.CategoryId, &transaction.Type); err != nil {
			log.Printf("Error scanning transaction row: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, transaction)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction

	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		log.Printf("Error decoding transaction: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse the date to ensure it's in the correct format
	date, err := time.Parse("02/01/2006", transaction.Date) // Assuming the input format is DD/MM/YYYY
	if err != nil {
		log.Printf("Error parsing date: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	transaction.Date = date.Format("2006-01-02") // Format date to YYYY-MM-DD for MySQL

	// Insert the new transaction into the database
	result, err := h.DB.Exec("INSERT INTO transactions (description, amount, date, budget_id, category_id, type, user_id) VALUES (?, ?, ?, ?, ?, ?, 1)",
		transaction.Description, transaction.Amount, transaction.Date, transaction.BudgetId, transaction.CategoryId, transaction.Type)
	if err != nil {
		log.Printf("Error inserting transaction: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Retrieve the ID of the newly inserted transaction
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the ID of the transaction and return the created transaction
	transaction.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}
