package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/arnavmaiti/oauth-microservice/internal/db"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	clientID := query.Get("client_id")
	redirectURI := query.Get("redirect_uri")
	scope := query.Get("scope")
	userID := query.Get("user_id") // For testing only, in prod use session

	if clientID == "" || redirectURI == "" || userID == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Validate client exists
	var exists bool
	err := db.GetDB().QueryRow("SELECT EXISTS(SELECT 1 FROM oauth_clients WHERE client_id=$1)", clientID).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Invalid client_id", http.StatusBadRequest)
		return
	}

	// Generate authorization code
	code := uuid.New().String()
	expiresAt := time.Now().Add(5 * time.Minute) // short-lived code
	scopes := pq.StringArray{scope}

	_, err = db.GetDB().Exec(`
        INSERT INTO oauth_authorization_codes 
        (id, code, user_id, client_id, redirect_uri, scopes, expires_at, created_at)
        VALUES ($1,$2,$3,(SELECT id FROM oauth_clients WHERE client_id=$4),$5,$6,$7,$8)
    `, uuid.New(), code, userID, clientID, redirectURI, scopes, expiresAt, time.Now())
	if err != nil {
		http.Error(w, "Failed to create authorization code", http.StatusInternalServerError)
		return
	}

	// Redirect back to client with code
	redirect := fmt.Sprintf("%s?code=%s", redirectURI, code)
	http.Redirect(w, r, redirect, http.StatusFound)
}
