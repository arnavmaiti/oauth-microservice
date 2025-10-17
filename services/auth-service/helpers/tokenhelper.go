package helpers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/arnavmaiti/oauth-microservice/services/auth-service/models"
	"github.com/arnavmaiti/oauth-microservice/services/common/constants"
	"github.com/arnavmaiti/oauth-microservice/services/common/db"
	"github.com/arnavmaiti/oauth-microservice/services/common/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// validateClient checks client_id and client_secret and returns internal client UUID, allowed grant_types and scopes
func ValidateClient(clientID, clientSecret string) (clientUUID string, grantTypes []string, scopes string, err error) {
	var client models.OAuthClient

	row := db.Get().QueryRow(constants.CheckClient, clientID)
	if err := row.Scan(&client.ID, &client.ClientID, &client.ClientSecret, pq.Array(&client.RedirectURIs), &client.Scopes, pq.Array(&client.GrantTypes)); err != nil {
		return "", nil, "", fmt.Errorf("%s", string(errors.CLIENT_NOT_FOUND))
	}

	// Client secret check
	if clientSecret != client.ClientSecret {
		return "", nil, "", fmt.Errorf("%s", string(errors.INVALID_CLIENT_SECRET))
	}

	return client.ID.String(), client.GrantTypes, client.Scopes, nil
}

// generateAndPersistTokens creates JWT access token and refresh token, persists into oauth_tokens table
func generateAndPersistTokens(subjectUserID *uuid.UUID, clientUUID string, scope string, expiresIn time.Duration, issuer string, jwtSigningKey []byte) (models.TokenResponse, error) {
	now := time.Now().UTC()
	exp := now.Add(expiresIn)

	var subjectUserIDValue sql.NullString
	if subjectUserID == nil {
		subjectUserIDValue = sql.NullString{Valid: false}
	} else {
		subjectUserIDValue = sql.NullString{String: (*subjectUserID).String()}
	}

	// Create JWT access token
	claims := jwt.MapClaims{
		"sub":   subjectUserIDValue.String,
		"aud":   clientUUID,
		"scope": scope,
		"iat":   now.Unix(),
		"exp":   exp.Unix(),
		"iss":   issuer,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(jwtSigningKey)
	if err != nil {
		return models.TokenResponse{}, err
	}
	// Create refresh token (opaque UUID)
	refreshToken := uuid.NewString()
	// Persist to DB
	id := uuid.New()
	_, err = db.Get().Exec(constants.AddToken, id, subjectUserIDValue, clientUUID, accessToken, refreshToken, scope, exp, now)
	if err != nil {
		return models.TokenResponse{}, err
	}
	resp := models.TokenResponse{
		AccessToken: accessToken,
		TokenType:   "bearer",
		ExpiresIn:   int64(expiresIn.Seconds()),
		Scope:       scope,
	}
	if subjectUserID != nil {
		resp.RefreshToken = refreshToken
	}
	return resp, nil
}

// authorization_code grant
func HandleAuthorizationCodeGrant(w http.ResponseWriter, r *http.Request, clientUUID string, clientGrantTypes []string, clientScopes string, issuer string, jwtSigningKey []byte) {
	code := r.FormValue("code")
	redirectURI := r.FormValue("redirect_uri")
	if code == "" || redirectURI == "" {
		errors.OAuthError(w, http.StatusBadRequest, errors.INVALID_REQUEST, errors.CODE_REDIRECT_URI_REQUIRED)
		return
	}

	// Validate code
	var authCode models.OAuthAuthorizationCode
	err := db.Get().QueryRow(constants.GetAuthCode, code).Scan(&authCode.UserID, &authCode.ClientID, &authCode.Scopes, &authCode.ExpiresAt, &authCode.RedirectURI)
	if err != nil {
		errors.OAuthError(w, http.StatusBadRequest, errors.INVALID_GRANT, errors.AUTHORIZATION_CODE_INVALID)
		return
	}
	if time.Now().After(authCode.ExpiresAt) {
		errors.OAuthError(w, http.StatusBadRequest, errors.INVALID_GRANT, errors.AUTHORIZATION_CODE_EXPIRED)
		return
	}
	// Ensure the authorization code belongs to this client
	if authCode.ClientID.String() != clientUUID {
		errors.OAuthError(w, http.StatusBadRequest, errors.INVALID_GRANT, errors.AUTHORIZATION_CODE_NOT_ISSUED)
		return
	}
	// Verify redirect_uri matches the one stored with code
	if authCode.RedirectURI != redirectURI {
		errors.OAuthError(w, http.StatusBadRequest, errors.INVALID_GRANT, errors.REDIRECT_URI_NO_MATCH)
		return
	}
	// Delete authorization code after use (one-time use)
	_, _ = db.Get().Exec(constants.DeleteAuthCode, code)

	// Use scopes from code if present, otherwise client's default
	scopeVal := ""
	if authCode.Scopes != "" {
		scopeVal = authCode.Scopes
	} else {
		scopeVal = clientScopes
	}
	// Generate tokens and persist
	resp, err := generateAndPersistTokens(&authCode.UserID, clientUUID, scopeVal, time.Hour, issuer, jwtSigningKey)
	if err != nil {
		errors.OAuthError(w, http.StatusInternalServerError, errors.SERVER_ERROR, errors.FAILED_TO_GENERATE_TOKENS)
		return
	}

	w.Header().Set(constants.ContentType, constants.ContentJSON)
	json.NewEncoder(w).Encode(resp)
}

// refresh_token grant
func HandleRefreshTokenGrant(w http.ResponseWriter, r *http.Request, clientUUID string, clientScopes string, issuer string, jwtSigningKey []byte) {
	refreshToken := r.FormValue("refresh_token")
	if refreshToken == "" {
		errors.OAuthError(w, http.StatusBadRequest, errors.INVALID_REQUEST, errors.REFRESH_TOKEN_REQUIRED)
		return
	}
	// Validate refresh token and find associated user & client
	var token models.OAuthToken
	err := db.Get().QueryRow(constants.GetToken, refreshToken).Scan(&token.UserID, &token.ClientID, &token.Scopes, &token.ExpiresAt)
	if err != nil {
		errors.OAuthError(w, http.StatusBadRequest, errors.INVALID_GRANT, errors.REFRESH_TOKEN_INVALID)
		return
	}
	// Check expiration of access token or refresh token policy
	if time.Now().After(token.ExpiresAt.Add(24 * time.Hour)) {
		errors.OAuthError(w, http.StatusBadRequest, errors.INVALID_GRANT, errors.REFRESH_TOKEN_EXPIRED)
		return
	}
	// Ensure token belongs to client
	if token.ClientID.String() != clientUUID {
		errors.OAuthError(w, http.StatusBadRequest, errors.INVALID_GRANT, errors.REFRESH_TOKEN_NOT_ISSUED_BY_CLIENT)
		return
	}
	// Rotate refresh token: delete old token row and create new
	_, _ = db.Get().Exec(constants.DeleteToken, refreshToken)
	// Use scopes from token if present, otherwise client's default
	scopeVal := ""
	if token.Scopes != "" {
		scopeVal = token.Scopes
	} else {
		scopeVal = clientScopes
	}

	resp, err := generateAndPersistTokens(&token.UserID, clientUUID, scopeVal, time.Hour, issuer, jwtSigningKey)
	if err != nil {
		errors.OAuthError(w, http.StatusInternalServerError, errors.SERVER_ERROR, errors.FAILED_TO_GENERATE_TOKENS)
		return
	}
	w.Header().Set(constants.ContentType, constants.ContentJSON)
	json.NewEncoder(w).Encode(resp)
}

// client_credentials grant
func HandleClientCredentialsGrant(w http.ResponseWriter, r *http.Request, clientUUID string, clientGrantTypes []string, clientScopes string, issuer string, jwtSigningKey []byte) {
	// Check client is authorized for client_credentials
	allowed := false
	for _, g := range clientGrantTypes {
		if g == "client_credentials" {
			allowed = true
			break
		}
	}
	if !allowed {
		errors.OAuthError(w, http.StatusUnauthorized, errors.UNAUTHORIZED_CLIENT, errors.CLIENT_NOT_ALLOWED_TO_USE_CLIENT_CRED)
		return
	}
	// Here subject is the client itself
	resp, err := generateAndPersistTokens(nil, clientUUID, clientScopes, time.Hour, issuer, jwtSigningKey)
	if err != nil {
		errors.OAuthError(w, http.StatusInternalServerError, errors.SERVER_ERROR, errors.FAILED_TO_GENERATE_TOKENS)
		return
	}
	w.Header().Set(constants.ContentType, constants.ContentJSON)
	json.NewEncoder(w).Encode(resp)
}

// password grant
func HandlePasswordGrant(w http.ResponseWriter, r *http.Request, clientUUID string, clientGrantTypes []string, clientScopes string, issuer string, jwtSigningKey []byte) {
	allowed := false
	for _, g := range clientGrantTypes {
		if g == "password" {
			allowed = true
			break
		}
	}
	if !allowed {
		errors.OAuthError(w, http.StatusUnauthorized, errors.UNAUTHORIZED_CLIENT, errors.CLIENT_NOT_ALLOWED_TO_USE_PASSWORD)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		errors.OAuthError(w, http.StatusBadRequest, errors.INVALID_REQUEST, errors.USERNAME_PASSWORD_REQUIRED)
		return
	}
	// Fetch user from user-service
	user, err := FetchUserFromUserService(username)
	if err != nil {
		errors.OAuthError(w, http.StatusUnauthorized, errors.INVALID_GRANT, errors.INVALID_CREDENTIALS)
		return
	}
	// Verify password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		errors.OAuthError(w, http.StatusUnauthorized, errors.INVALID_GRANT, errors.INVALID_CREDENTIALS)
		return
	}
	resp, err := generateAndPersistTokens(&user.ID, clientUUID, clientScopes, time.Hour, issuer, jwtSigningKey)
	if err != nil {
		errors.OAuthError(w, http.StatusInternalServerError, errors.SERVER_ERROR, errors.FAILED_TO_GENERATE_TOKENS)
		return
	}
	w.Header().Set(constants.ContentType, constants.ContentJSON)
	json.NewEncoder(w).Encode(resp)
}
