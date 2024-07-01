package handlers

import (
	"budgeting-app/internal/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type CategoryHandler struct {
	DB *sql.DB
}

func NewCategoryHandler(db *sql.DB) *CategoryHandler {
	return &CategoryHandler{DB: db}
}

func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	var categories []models.Category

	rows, err := h.DB.Query("SELECT id, description, amount FROM category")
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			log.Printf("Error scanning categories row: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
