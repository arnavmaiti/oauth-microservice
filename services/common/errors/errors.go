package errors

import (
	"encoding/json"
	"net/http"

	"github.com/arnavmaiti/oauth-microservice/services/common/models"
)

const (
	// Users
	UserNotFound string = "User not found"
)

type ErrorCode string

const (
	INVALID_CLIENT         ErrorCode = "invalid_client"
	INVALID_GRANT          ErrorCode = "invalid_grant"
	INVALID_REQUEST        ErrorCode = "invalid_request"
	INVALID_SCOPE          ErrorCode = "invalid_scope"
	SERVER_ERROR           ErrorCode = "server_error"
	UNAUTHORIZED_CLIENT    ErrorCode = "unauthorized_client"
	UNSUPPORTED_GRANT_TYPE ErrorCode = "unsupported_grant_type"
)

type ErrorMessage string

const (
	AUTHORIZATION_CODE_EXPIRED             ErrorMessage = "authorization code expired"
	AUTHORIZATION_CODE_INVALID             ErrorMessage = "authorization code invalid"
	AUTHORIZATION_CODE_NOT_ISSUED          ErrorMessage = "authorization code was not issued to this client"
	CLIENT_AUTHENTICATION_FAILED           ErrorMessage = "client authentication failed"
	CLIENT_NOT_FOUND                       ErrorMessage = "client not found"
	CODE_REDIRECT_URI_REQUIRED             ErrorMessage = "code and redirect_uri required"
	FAILED_TO_GENERATE_TOKENS              ErrorMessage = "failed to generate tokens"
	GRANT_TYPE_NOT_SUPPORTED               ErrorMessage = "grant type is not supported"
	GRANT_TYPE_REQUIRED                    ErrorMessage = "grant type is required"
	INVALID_CLIENT_SECRET                  ErrorMessage = "invalid client secret"
	INVALID_FORM_DATA                      ErrorMessage = "invalid form data"
	REDIRECT_URI_NO_MATCH                  ErrorMessage = "redirect URI does not match"
	MESSAGE_ONLY_POST_METHOD_ALLOWED       ErrorMessage = "request method must be POST"
	MESSAGE_ONLY_JSON_TYPE_CONTENT_ALLOWED ErrorMessage = "content-Type must be application/x-www-form-urlencoded"
)

// OAuthError writes a RFC-compliant error response
func OAuthError(w http.ResponseWriter, status int, errCode ErrorCode, description ErrorMessage) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := models.ErrorResponse{
		Error:            string(errCode),
		ErrorDescription: string(description),
	}
	_ = json.NewEncoder(w).Encode(resp)
}
