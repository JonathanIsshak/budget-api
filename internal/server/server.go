package server

import (
	"budgeting-app/internal/handlers"
	"budgeting-app/internal/middleware"
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
	transactionHandler := handlers.NewTransactionHandler(s.DB)
	categoryHandler := handlers.NewCategoryHandler(s.DB)
	userHandler := handlers.NewUserHandler(s.DB)

	s.Router.Handle("/api/budgets", middleware.JWTAuth(http.HandlerFunc(budgetHandler.GetBudgets))).Methods("GET")
	s.Router.HandleFunc("/api/budgets", budgetHandler.CreateBudget).Methods("POST")

	s.Router.HandleFunc("/api/transactions", transactionHandler.GetTransactions).Methods("GET")
	s.Router.HandleFunc("/api/transactions", transactionHandler.CreateTransaction).Methods("POST")
	s.Router.HandleFunc("/api/categories", categoryHandler.GetCategories).Methods("GET")
	s.Router.HandleFunc("/api/categories", categoryHandler.CreateCategory).Methods("POST")
	s.Router.HandleFunc("/api/users/{id}", userHandler.GetUser).Methods("GET")
	s.Router.HandleFunc("/api/users", userHandler.CreateUser).Methods("POST")
	s.Router.HandleFunc("/api/users/login", userHandler.LoginUser).Methods("POST")
}

func (s *Server) RouterFunc() http.Handler {
	return s.Router
}
