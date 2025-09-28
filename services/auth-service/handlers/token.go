package handlers

import (
	"net/http"
	"strings"

	"github.com/arnavmaiti/oauth-microservice/services/auth-service/helpers"
	"github.com/arnavmaiti/oauth-microservice/services/common/constants"
	"github.com/arnavmaiti/oauth-microservice/services/common/errors"
)

var jwtSigningKey = []byte("replace-with-secure-key")
var issuer = "https://auth.example.local"

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	// Must be POST with form encoding
	if r.Method != http.MethodPost {
		errors.OAuthError(w, http.StatusMethodNotAllowed, errors.INVALID_REQUEST, errors.MESSAGE_ONLY_POST_METHOD_ALLOWED)
		return
	}
	// Must be URL form encoded
	if ct := r.Header.Get(constants.ContentType); !strings.HasPrefix(ct, constants.URLFormEncoded) {
		errors.OAuthError(w, http.StatusBadRequest, errors.INVALID_REQUEST, errors.MESSAGE_ONLY_JSON_TYPE_CONTENT_ALLOWED)
		return
	}
	if err := r.ParseForm(); err != nil {
		errors.OAuthError(w, http.StatusBadRequest, errors.INVALID_REQUEST, errors.INVALID_FORM_DATA)
		return
	}
	// Validate grant type is present
	grantType := r.FormValue("grant_type")
	if grantType == "" {
		errors.OAuthError(w, http.StatusBadRequest, errors.INVALID_REQUEST, errors.GRANT_TYPE_REQUIRED)
		return
	}
	// Authenticate client (Basic Auth preferred)
	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		clientID = r.FormValue("client_id")
		clientSecret = r.FormValue("client_secret")
	}
	if clientID == "" {
		// invalid_client per spec
		w.Header().Set(constants.WWWAuthenticate, constants.OAuth2BasicRealm)
		errors.OAuthError(w, http.StatusUnauthorized, errors.INVALID_CLIENT, errors.CLIENT_AUTHENTICATION_FAILED)
		return
	}
	// Load client from DB and validate secret
	clientDBID, clientGrantTypes, clientScopes, err := helpers.ValidateClient(clientID, clientSecret)
	if err != nil {
		w.Header().Set(constants.WWWAuthenticate, constants.OAuth2BasicRealm)
		errors.OAuthError(w, http.StatusUnauthorized, errors.INVALID_CLIENT, errors.CLIENT_AUTHENTICATION_FAILED)
		return
	}

	// Switch grant type
	switch grantType {
	case "authorization_code":
		helpers.HandleAuthorizationCodeGrant(w, r, clientDBID, clientGrantTypes, clientScopes, issuer, jwtSigningKey)
	case "refresh_token":
		helpers.HandleRefreshTokenGrant(w, r, clientDBID, clientScopes, issuer, jwtSigningKey)
	case "client_credentials":
	case "password":
	default:
		errors.OAuthError(w, http.StatusBadRequest, errors.UNSUPPORTED_GRANT_TYPE, errors.GRANT_TYPE_NOT_SUPPORTED)
		return
	}
}
