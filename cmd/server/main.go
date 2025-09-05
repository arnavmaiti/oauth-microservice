package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/arnavmaiti/oauth-microservice/internal/handlers"
)

var db *sql.DB

func main() {
	// Check PostGRES connection
	pgUser := os.Getenv("POSTGRES_USER")
	if pgUser == "" {
		log.Fatal("POSTGRES_USER is not set")
	}
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	if pgPassword == "" {
		log.Fatal("POSTGRES_PASSWORD is not set")
	}
	pgHost := os.Getenv("POSTGRES_HOST")
	if pgHost == "" {
		log.Fatal("POSTGRES_HOST is not set")
	}
	pgPort := os.Getenv("POSTGRES_PORT")
	if pgPort == "" {
		log.Fatal("POSTGRES_PORT is not set")
	}
	pgDatabase := os.Getenv("POSTGRES_DB")
	if pgDatabase == "" {
		log.Fatal("POSTGRES_DB is not set")
	}
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", pgUser, pgPassword, pgHost, pgPort, pgDatabase)

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Unable to find database: %v", err)
	}

	// Ping and check
	if err := db.Ping(); err != nil {
		log.Fatalf("Unable to connect to database with error: %v", err)
	}
	log.Println("Successfully connected to database")

	healthContext := &handlers.HealthContext{}
	healthContext.SetDatabase(db)

	mux := http.NewServeMux()

	// Register endpoints
	mux.HandleFunc("/health", healthContext.HealthCheck)
	mux.HandleFunc("/ready", healthContext.ReadyCheck)

	port := ":8080"
	fmt.Printf("OAuth server running on %s\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
