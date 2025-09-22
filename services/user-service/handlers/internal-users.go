package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/arnavmaiti/oauth-microservice/services/common/constants"
	db "github.com/arnavmaiti/oauth-microservice/services/common/db"
	"github.com/arnavmaiti/oauth-microservice/services/user-service/models"
)

func HandleInternalUserByID(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	switch r.Method {
	case http.MethodGet:
		var u models.User
		err := db.Get().QueryRow(constants.GetInternalUser, username).
			Scan(&u.ID, &u.Username, &u.PasswordHash)
		if err != nil {
			http.Error(w, constants.UserNotFound, http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(u)
	}
}
