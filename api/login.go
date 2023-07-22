package api

import (
	"encoding/json"
	"fmt"
	"github.com/Krabik6/meal-shedule-telegram-bot/internal/model"
	"io"
	"net/http"
	"strings"
)

// Login that takes 4 arguments: bot, update, client, user, and returns error. It is used to log in the user using the API.
func (a *Api) Login(user model.LoginCredentials) (string, error) {
	loginCredentials := model.LoginCredentials{
		Username: user.Username,
		Password: user.Password,
	}

	requestBody, err := json.Marshal(loginCredentials)
	if err != nil {
		return "", err
	}

	resp, err := a.Client.Post("http://localhost:8000/auth/sign-in", "application/json", strings.NewReader(string(requestBody)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("неверный логин или пароль")
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("код состояния %d. \nответ: %v", resp.StatusCode, string(body))
	}

	var authResponse model.AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	if err != nil {
		return "", err
	}

	return authResponse.Token, nil
}
