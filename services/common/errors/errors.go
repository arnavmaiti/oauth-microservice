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
	INVALID_CLIENT            ErrorCode = "invalid_client"
	INVALID_GRANT             ErrorCode = "invalid_grant"
	INVALID_REQUEST           ErrorCode = "invalid_request"
	INVALID_SCOPE             ErrorCode = "invalid_scope"
	SERVER_ERROR              ErrorCode = "server_error"
	UNAUTHORIZED_CLIENT       ErrorCode = "unauthorized_client"
	UNSUPPORTED_GRANT_TYPE    ErrorCode = "unsupported_grant_type"
	UNSUPPORTED_RESPONSE_TYPE ErrorCode = "unsupported_response_type"
)

type ErrorMessage string

const (
	AUTHORIZATION_CODE_CREATION_FAILED     ErrorMessage = "authorization code creation failed"
	AUTHORIZATION_CODE_EXPIRED             ErrorMessage = "authorization code expired"
	AUTHORIZATION_CODE_INVALID             ErrorMessage = "authorization code invalid"
	AUTHORIZATION_CODE_NOT_ISSUED          ErrorMessage = "authorization code was not issued to this client"
	CLIENT_AUTHENTICATION_FAILED           ErrorMessage = "client authentication failed"
	CLIENT_ID_REDIRECT_URI_REQUIRED        ErrorMessage = "client_id and redirect_uri required"
	CLIENT_NOT_ALLOWED_TO_USE_CLIENT_CRED  ErrorMessage = "client not allowed to use client_credentials"
	CLIENT_NOT_ALLOWED_TO_USE_PASSWORD     ErrorMessage = "client not allowed to use password grant"
	CLIENT_NOT_FOUND                       ErrorMessage = "client not found"
	CODE_REDIRECT_URI_REQUIRED             ErrorMessage = "code and redirect_uri required"
	FAILED_TO_GENERATE_TOKENS              ErrorMessage = "failed to generate tokens"
	GRANT_TYPE_NOT_SUPPORTED               ErrorMessage = "grant type is not supported"
	GRANT_TYPE_REQUIRED                    ErrorMessage = "grant type is required"
	INVALID_CLIENT_SECRET                  ErrorMessage = "invalid client secret"
	INVALID_CREDENTIALS                    ErrorMessage = "invalid credentials"
	INVALID_FORM_DATA                      ErrorMessage = "invalid form data"
	INVALID_REDIRECT_URI                   ErrorMessage = "invalid redirect URI"
	MESSAGE_ONLY_POST_METHOD_ALLOWED       ErrorMessage = "request method must be POST"
	MESSAGE_ONLY_JSON_TYPE_CONTENT_ALLOWED ErrorMessage = "content-Type must be application/x-www-form-urlencoded"
	REDIRECT_URI_NO_MATCH                  ErrorMessage = "redirect URI does not match"
	REFRESH_TOKEN_EXPIRED                  ErrorMessage = "refresh token expired"
	REFRESH_TOKEN_INVALID                  ErrorMessage = "refresh token invalid"
	REFRESH_TOKEN_NOT_ISSUED_BY_CLIENT     ErrorMessage = "refresh token was not issued to this client"
	REFRESH_TOKEN_REQUIRED                 ErrorMessage = "refresh_token required"
	RESPONSE_TYPE_NOT_SUPPORTED            ErrorMessage = "response type not supported"
	USERNAME_PASSWORD_REQUIRED             ErrorMessage = "username and password required"
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
