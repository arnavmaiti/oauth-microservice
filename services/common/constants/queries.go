package constants

const (
	DeleteUser      string = `DELETE FROM users WHERE id=$1`
	CreateUser      string = `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id`
	GetUser         string = `SELECT id, username, email, created_at, updated_at FROM users WHERE id=$1`
	GetInternalUser string = `SELECT id, username, password_hash FROM users WHERE username=$1`
	GetUsers        string = `SELECT id, username, email, created_at, updated_at FROM users`
)
