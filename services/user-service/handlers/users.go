package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/arnavmaiti/oauth-microservice/services/common/constants"
	db "github.com/arnavmaiti/oauth-microservice/services/common/db"
	"github.com/arnavmaiti/oauth-microservice/services/user-service/models"
	"golang.org/x/crypto/bcrypt"
)

// TODO: Need better error handling
func HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		var id string
		err := db.Get().QueryRow(
			constants.CreateUser,
			req.Username, req.Email, string(hash),
		).Scan(&id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id})

	case http.MethodGet:
		rows, err := db.Get().Query(constants.GetUsers)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []models.User
		for rows.Next() {
			var u models.User
			if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			users = append(users, u)
		}
		json.NewEncoder(w).Encode(users)
	}
}

func HandleUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	switch r.Method {
	case http.MethodGet:
		var u models.User
		err := db.Get().QueryRow(constants.GetUser, id).
			Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			http.Error(w, constants.UserNotFound, http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(u)

	case http.MethodDelete:
		_, err := db.Get().Exec(constants.DeleteUser, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
