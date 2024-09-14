// internal/database/supabase.go

package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	SupabaseURL string
	SupabaseKey string
)

func InitSupabase(url, key string) {
	SupabaseURL = url
	SupabaseKey = key
	log.Printf("Supabase initialized with URL: %s", SupabaseURL)
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ID           string    `json:"id"`
		Aud          string    `json:"aud"`
		Role         string    `json:"role"`
		Email        string    `json:"email"`
		ConfirmedAt  time.Time `json:"confirmed_at"`
		LastSignInAt time.Time `json:"last_sign_in_at"`
		AppMetadata  struct {
			Provider string `json:"provider"`
		} `json:"app_metadata"`
		UserMetadata struct{}  `json:"user_metadata"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	} `json:"user"`
}

func RegisterUser(email, password string) (AuthResponse, error) {
	url := fmt.Sprintf("%s/auth/v1/signup", SupabaseURL)

	userData := map[string]string{
		"email":    email,
		"password": password,
	}

	jsonData, err := json.Marshal(userData)
	if err != nil {
		log.Printf("Error marshaling user data: %v", err)
		return AuthResponse{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return AuthResponse{}, err
	}

	req.Header.Set("apikey", SupabaseKey)
	req.Header.Set("Authorization", "Bearer "+SupabaseKey)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("Sending registration request to: %s", url)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request to Supabase Auth: %v", err)
		return AuthResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return AuthResponse{}, err
	}

	log.Printf("Response body: %s", string(body))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("Supabase Auth returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
		return AuthResponse{}, fmt.Errorf("failed to register user: %s", body)
	}

	var authResp AuthResponse
	err = json.Unmarshal(body, &authResp)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return AuthResponse{}, err
	}

	log.Printf("User registered successfully: %s", authResp.User.Email)
	return authResp, nil
}
