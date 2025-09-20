package main

import (
	"log"
	"net/http"

	"github.com/arnavmaiti/oauth-microservice/services/common/constants"
	"github.com/arnavmaiti/oauth-microservice/services/user-service/handlers"
	_ "github.com/lib/pq"
)

func main() {
	mux := http.NewServeMux()
	// User APIs
	mux.HandleFunc("/users", handlers.HandleUsers)
	mux.HandleFunc("/users/{id}", handlers.HandleUserByID)
	// Internal User APIs
	mux.HandleFunc("/internal/users/{username}", handlers.HandleInternalUserByID)

	log.Printf("User Service running on %s\n", constants.UserPort)
	log.Fatal(http.ListenAndServe(constants.UserPort, mux))
}
