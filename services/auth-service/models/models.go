package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the OAuth2 system
type User struct {
	ID           uuid.UUID `db:"id"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	Email        string    `db:"email"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// OAuthClient represents an OAuth2 client application
type OAuthClient struct {
	ClientID     string    `db:"client_id"`
	ClientSecret string    `db:"client_secret"`
	RedirectURIs []string  `db:"redirect_uris"`
	Scopes       string    `db:"scopes"`
	GrantTypes   []string  `db:"grant_types"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// OAuthToken represents an access/refresh token pair
type OAuthToken struct {
	ID           uuid.UUID `db:"id"`
	UserID       uuid.UUID `db:"user_id"`
	ClientID     uuid.UUID `db:"client_id"`
	AccessToken  string    `db:"access_token"`
	RefreshToken string    `db:"refresh_token"`
	Scopes       []string  `db:"scopes"`
	ExpiresAt    time.Time `db:"expires_at"`
	CreatedAt    time.Time `db:"created_at"`
}

// OAuthAuthorizationCode represents an auth code in the Authorization Code flow
type OAuthAuthorizationCode struct {
	ID          uuid.UUID `db:"id"`
	Code        string    `db:"code"`
	UserID      uuid.UUID `db:"user_id"`
	ClientID    uuid.UUID `db:"client_id"`
	RedirectURI string    `db:"redirect_uri"`
	Scopes      []string  `db:"scopes"`
	ExpiresAt   time.Time `db:"expires_at"`
	CreatedAt   time.Time `db:"created_at"`
}
