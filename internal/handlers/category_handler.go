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

	rows, err := h.DB.Query("SELECT id, name FROM categories")
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

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		log.Printf("Error decoding category: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Insert the new category into the database
	result, err := h.DB.Exec("INSERT INTO categories (name, user_id) VALUES (?, 1)", category.Name)
	if err != nil {
		log.Printf("Error inserting category: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Retrieve the ID of the newly inserted category
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the ID of the category and return the created category
	category.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}
