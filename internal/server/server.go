// internal/server/server.go

package server

import (
	"budgeting-app/internal/handlers"
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	DB     *sql.DB
}

func NewServer(db *sql.DB) *Server {
	s := &Server{
		Router: mux.NewRouter(),
		DB:     db,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	budgetHandler := handlers.NewBudgetHandler(s.DB)

	s.Router.HandleFunc("/api/budgets", budgetHandler.GetBudgets).Methods("GET")
	// Define more routes here for other handlers
}

func (s *Server) RouterFunc() http.Handler {
	return s.Router
}
