package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/arnavmaiti/oauth-microservice/internal/db"
)

type HealthResponse struct {
	Status string `json:"status"`
}

// Health check endpoint
func HealthCheck(w http.ResponseWriter, r *http.Request) {

	resp := HealthResponse{Status: "OK"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func ReadyCheck(w http.ResponseWriter, r *http.Request) {

	dbConn := db.GetDB()

	var exists bool
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = 'schema_migrations'
		);
	`

	if err := dbConn.QueryRow(query).Scan(&exists); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Error checking schema: %v\n", err)
		return
	}

	if !exists {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, "Migrations not applied yet")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "READY")
}
