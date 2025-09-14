package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/arnavmaiti/oauth-microservice/internal/handlers"
)

var db *sql.DB

func main() {

	mux := http.NewServeMux()

	// Register endpoints
	mux.HandleFunc("/health", handlers.HealthCheck)
	mux.HandleFunc("/ready", handlers.ReadyCheck)
	mux.HandleFunc("/register", handlers.RegisterHandler)

	port := ":8080"
	fmt.Printf("OAuth server running on %s\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
