package handlers

import (
	"encoding/json"
	"net/http"
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
