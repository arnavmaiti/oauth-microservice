package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/arnavmaiti/oauth-microservice/internal/handlers"
)

func main() {

	mux := http.NewServeMux()

	// Register endpoints
	mux.HandleFunc("/health", handlers.HealthCheck)
	mux.HandleFunc("/ready", handlers.ReadyCheck)
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/authorize", handlers.AuthorizeHandler)
	mux.HandleFunc("/token", handlers.TokenHandler)
	mux.HandleFunc("/introspect", handlers.IntrospectHandler)

	port := ":8080"
	fmt.Printf("OAuth server running on %s\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
