package handlers

import (
	"net/http"
	"net/url"
	"time"

	"github.com/arnavmaiti/oauth-microservice/services/auth-service/models"
	"github.com/arnavmaiti/oauth-microservice/services/common/constants"
	"github.com/arnavmaiti/oauth-microservice/services/common/db"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	responseType := query.Get("response_type")
	clientID := query.Get("client_id")
	redirectURI := query.Get("redirect_uri")
	scopes := query.Get("scopes")
	state := query.Get("state")

	userID := query.Get("user_id") // TODO: Fix this once user ID can be fetched with session

	if responseType != "code" {
		http.Error(w, "unsupported_response_type", http.StatusBadRequest)
		return
	}

	if clientID == "" || redirectURI == "" {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	// Get the client
	var client models.OAuthClient
	err := db.Get().QueryRow(constants.CheckClient, clientID).Scan(&client.ID, &client.ClientID, &client.ClientSecret, pq.Array(&client.RedirectURIs), &client.Scopes, pq.Array(&client.GrantTypes))
	if err != nil {
		http.Error(w, "unauthorized_client", http.StatusUnauthorized)
		return
	}
	// Check redirect URI
	validRedirect := false
	for _, uri := range client.RedirectURIs {
		if uri == redirectURI {
			validRedirect = true
			break
		}
	}
	if !validRedirect {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	// Generate authorization code
	code := uuid.New().String()
	expiresAt := time.Now().Add(5 * time.Minute) // short-lived code

	_, err = db.Get().Exec(constants.AddAuthCode, code, userID, client.ID, redirectURI, scopes, expiresAt)
	if err != nil {
		http.Error(w, "server_error", http.StatusInternalServerError)
		return
	}

	// Redirect back to client with code and state
	redirect, err := url.Parse(redirectURI)
	if err != nil {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}
	values := redirect.Query()
	values.Set("code", code)
	if state != "" {
		values.Set("state", state)
	}
	redirect.RawQuery = values.Encode()
	http.Redirect(w, r, redirect.String(), http.StatusFound)
}
