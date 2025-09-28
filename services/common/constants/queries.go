package constants

const (
	// User queries
	DeleteUser      string = `DELETE FROM users WHERE id=$1`
	CreateUser      string = `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id`
	GetUser         string = `SELECT id, username, email, created_at, updated_at FROM users WHERE id=$1`
	GetInternalUser string = `SELECT id, username, password_hash FROM users WHERE username=$1`
	GetUsers        string = `SELECT id, username, email, created_at, updated_at FROM users`

	// Auth queries
	CheckClient    string = `SELECT id, client_id, client_secret, redirect_uris, scopes, grant_types FROM oauth_clients WHERE client_id = $1`
	AddAuthCode    string = `INSERT INTO oauth_authorization_codes (code, user_id, client_id, redirect_uri, scopes, expires_at) VALUES ($1,$2,$3,$4,$5,$6)`
	GetAuthCode    string = `SELECT user_id, client_id, scopes, expires_at, redirect_uri FROM oauth_authorization_codes WHERE code = $1`
	DeleteAuthCode string = `DELETE FROM oauth_authorization_codes WHERE code=$1`
	AddToken       string = `INSERT INTO oauth_tokens (id, user_id, client_id, access_token, refresh_token, scopes, expires_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	GetToken       string = `SELECT user_id, client_id, scopes, expires_at FROM oauth_tokens WHERE refresh_token=$1`
	DeleteToken    string = `DELETE FROM oauth_tokens WHERE refresh_token=$1`
)
