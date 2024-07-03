package handlers

import (
	"budgeting-app/internal/auth"
	"budgeting-app/internal/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	DB *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{DB: db}
}

// GetUser fetches a user by ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	id := r.URL.Query().Get("id")

	row := h.DB.QueryRow("SELECT id, username, email, created_at, updated_at FROM users WHERE id = ?", id)

	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		log.Printf("Error scanning user row: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// return the user as json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Error decoding user: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// hash the password
	hashedPassword, err := hashPassword(user.PasswordHash)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.PasswordHash = hashedPassword

	_, err = h.DB.Exec(
		"INSERT INTO users (username, password_hash, email, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())",
		user.Username, user.PasswordHash, user.Email,
	)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// LoginUser authenticates a user and returns a JWT token
func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var creds models.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		log.Printf("Error decoding login credentials: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var storedCreds models.User
	err := h.DB.QueryRow("SELECT password_hash FROM users WHERE username = ?", creds.Username).Scan(&storedCreds.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		log.Printf("Error fetching stored credentials: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Compare the stored hashed password with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(storedCreds.PasswordHash), []byte(creds.PasswordHash)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(creds.Username)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return the token as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
