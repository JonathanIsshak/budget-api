package handlers

import (
	"budgeting-app/internal/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type TransactionHandler struct {
	DB *sql.DB
}

func NewTransactionHandler(db *sql.DB) *TransactionHandler {
	return &TransactionHandler{DB: db}
}

func (h *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	var transactions []models.Transaction

	rows, err := h.DB.Query("SELECT id, description, amount FROM transactions")
	if err != nil {
		log.Printf("Error fetching transactions: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.Description, &transaction.Amount); err != nil {
			log.Printf("Error scanning transaction row: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, transaction)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}
