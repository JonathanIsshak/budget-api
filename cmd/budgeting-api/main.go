package main

import (
	"log"
	"net/http"

	"budgeting-app/db"
	"budgeting-app/internal/server"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Connect to database
	db, err := db.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Initialize HTTP server
	srv := server.NewServer(db)
	addr := ":8888"

	// Start server
	log.Printf("Server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, srv.RouterFunc()))
}
