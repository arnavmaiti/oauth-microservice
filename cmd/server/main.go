package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/arnavmaiti/oauth-microservice/internal/handlers"
)

func main() {
	mux := http.NewServeMux()

	// Register endpoints
	mux.HandleFunc("/health", handlers.HealthCheck)

	port := ":8080"
	fmt.Printf("OAuth server running on %s\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
