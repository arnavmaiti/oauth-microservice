package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/arnavmaiti/oauth-microservice/internal/db"
	"github.com/arnavmaiti/oauth-microservice/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Username == "" || req.Email == "" || len(req.Password) < 8 {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user := models.User{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashed),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Insert into Postgres using prepared statement
	dbConn := db.GetDB()
	query := `
        INSERT INTO users (id, username, password_hash, email, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err = dbConn.Exec(query,
		user.ID,
		user.Username,
		user.PasswordHash,
		user.Email,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	resp := RegisterResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
