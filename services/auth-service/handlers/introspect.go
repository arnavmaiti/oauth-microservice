package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/arnavmaiti/oauth-microservice/services/common/db"
)

func IntrospectHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	token := r.PostForm.Get("token")
	clientID := r.PostForm.Get("client_id")
	clientSecret := r.PostForm.Get("client_secret")

	if token == "" || clientID == "" || clientSecret == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Validate client
	var dbClientID string
	err := db.Get().QueryRow(`
        SELECT client_id FROM oauth_clients WHERE client_id=$1 AND client_secret=$2
    `, clientID, clientSecret).Scan(&dbClientID)
	if err != nil {
		http.Error(w, "Invalid client credentials", http.StatusUnauthorized)
		return
	}

	// Validate token
	var userID string
	var expiresAt time.Time
	err = db.Get().QueryRow(`
        SELECT user_id, expires_at FROM oauth_tokens WHERE access_token=$1
    `, token).Scan(&userID, &expiresAt)
	if err != nil || time.Now().After(expiresAt) {
		json.NewEncoder(w).Encode(map[string]interface{}{"active": false})
		return
	}

	// Return token info
	json.NewEncoder(w).Encode(map[string]interface{}{
		"active":     true,
		"user_id":    userID,
		"expires_at": expiresAt,
	})
}
