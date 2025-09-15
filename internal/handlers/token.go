package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/arnavmaiti/oauth-microservice/internal/db"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	grantType := r.PostForm.Get("grant_type")

	switch grantType {
	case "authorization_code":
		code := r.PostForm.Get("code")
		clientID := r.PostForm.Get("client_id")
		clientSecret := r.PostForm.Get("client_secret")
		redirectURI := r.PostForm.Get("redirect_uri")

		if code == "" || clientID == "" || clientSecret == "" || redirectURI == "" {
			http.Error(w, "Missing required parameters", http.StatusBadRequest)
			return
		}

		// Validate client
		var dbClientID string
		err := db.GetDB().QueryRow("SELECT client_id FROM oauth_clients WHERE client_id=$1 AND client_secret=$2", clientID, clientSecret).Scan(&dbClientID)
		if err != nil {
			http.Error(w, "Invalid client credentials", http.StatusUnauthorized)
			return
		}

		// Validate auth code
		var userID string
		var expiresAt time.Time
		err = db.GetDB().QueryRow(`
			SELECT user_id, expires_at FROM oauth_authorization_codes 
			WHERE code=$1 AND redirect_uri=$2
		`, code, redirectURI).Scan(&userID, &expiresAt)
		if err != nil || time.Now().After(expiresAt) {
			http.Error(w, "Invalid or expired authorization code", http.StatusBadRequest)
			return
		}

		// Generate access & refresh tokens
		accessToken := uuid.New().String()
		refreshToken := uuid.New().String()
		tokenExpiry := time.Now().Add(1 * time.Hour)

		_, err = db.GetDB().Exec(`
			INSERT INTO oauth_tokens (id, user_id, client_id, access_token, refresh_token, scopes, expires_at, created_at)
			VALUES ($1, $2, (SELECT id FROM oauth_clients WHERE client_id=$3), $4, $5, $6, $7, $8)
		`, uuid.New(), userID, clientID, accessToken, refreshToken, pq.StringArray{}, tokenExpiry, time.Now())
		if err != nil {
			http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
			return
		}

		// Respond with JSON
		resp := TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    int64(time.Until(tokenExpiry).Seconds()),
			TokenType:    "Bearer",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	case "refresh_token":
		refreshToken := r.PostForm.Get("refresh_token")
		if refreshToken == "" {
			http.Error(w, "Missing refresh_token", http.StatusBadRequest)
			return
		}

		// Validate refresh token
		var userID string
		var clientDBID string
		err := db.GetDB().QueryRow(`
			SELECT user_id, client_id FROM oauth_tokens WHERE refresh_token=$1
		`, refreshToken).Scan(&userID, &clientDBID)
		if err != nil {
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
			return
		}

		// Generate new tokens
		newAccessToken := uuid.New().String()
		newRefreshToken := uuid.New().String()
		tokenExpiry := time.Now().Add(1 * time.Hour)

		_, err = db.GetDB().Exec(`
			INSERT INTO oauth_tokens (id, user_id, client_id, access_token, refresh_token, scopes, expires_at, created_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		`, uuid.New(), userID, clientDBID, newAccessToken, newRefreshToken, pq.StringArray{}, tokenExpiry, time.Now())
		if err != nil {
			http.Error(w, "Failed to generate new tokens", http.StatusInternalServerError)
			return
		}

		// Return new tokens
		resp := map[string]interface{}{
			"access_token":  newAccessToken,
			"refresh_token": newRefreshToken,
			"expires_in":    int64(time.Until(tokenExpiry).Seconds()),
			"token_type":    "Bearer",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	default:
		http.Error(w, "Unsupported grant_type", http.StatusBadRequest)
		return
	}
}
