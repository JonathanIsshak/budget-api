package db

import (
	"database/sql"
	"fmt"

	"budgeting-app/internal/config" // Import your config package

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Connect initializes a connection to the database.
func Connect() (*sql.DB, error) {
	cfg := config.LoadConfig() // Load configuration from config package

	dbSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	conn, err := sql.Open("mysql", dbSource)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	db = conn
	fmt.Println("Connected to database")
	return db, nil
}
