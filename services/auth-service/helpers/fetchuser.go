package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/arnavmaiti/oauth-microservice/services/auth-service/models"
)

func FetchUserFromUserService(username string) (*models.User, error) {
	url := fmt.Sprintf("http://user-service:8081/internal/users/%s", username) // k8s DNS
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user not found: %s", body)
	}

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
