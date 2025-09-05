package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

type HealthContext struct {
	db *sql.DB
}

func (context *HealthContext) SetDatabase(db *sql.DB) {
	context.db = db
}

// Health check endpoint
func (context *HealthContext) HealthCheck(w http.ResponseWriter, r *http.Request) {

	resp := HealthResponse{Status: "OK"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (context *HealthContext) ReadyCheck(w http.ResponseWriter, r *http.Request) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = 'schema_migrations'
		);
	`

	if err := context.db.QueryRow(query).Scan(&exists); err != nil {
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
