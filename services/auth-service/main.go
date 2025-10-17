package main

import (
	"log"
	"net/http"

	"github.com/arnavmaiti/oauth-microservice/services/auth-service/handlers"
	"github.com/arnavmaiti/oauth-microservice/services/common/constants"
	_ "github.com/lib/pq"
)

func main() {
	mux := http.NewServeMux()
	// Health APIs
	mux.HandleFunc("/health", handlers.HealthCheck)
	mux.HandleFunc("/ready", handlers.ReadyCheck)
	// Auth APIs
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/authorize", handlers.AuthorizeHandler)
	mux.HandleFunc("/token", handlers.TokenHandler)
	mux.HandleFunc("/introspect", handlers.IntrospectHandler)

	log.Printf("Auth Service running on %s\n", constants.AuthPort)
	log.Fatal(http.ListenAndServe(constants.AuthPort, mux))
}
