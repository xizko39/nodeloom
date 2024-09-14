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
}

type User struct {
	ID        string    `json:"id,omitempty"`
	Username  string    `json:"username"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func InsertUser(user User) (User, error) {
	url := fmt.Sprintf("%s/rest/v1/users", SupabaseURL)

	userData := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: user.Username,
		Password: user.Password,
	}

	jsonData, err := json.Marshal(userData)
	if err != nil {
		log.Printf("Error marshaling user data: %v", err)
		return User{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return User{}, err
	}

	req.Header.Set("apikey", SupabaseKey)
	req.Header.Set("Authorization", "Bearer "+SupabaseKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	log.Printf("Sending request to: %s", url)
	log.Printf("Request headers: Content-Type: %s, Prefer: %s",
		req.Header.Get("Content-Type"), req.Header.Get("Prefer"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request to Supabase: %v", err)
		return User{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return User{}, err
	}

	log.Printf("Response body: %s", string(body))

	if resp.StatusCode != http.StatusCreated {
		log.Printf("Supabase returned non-201 status: %d, body: %s", resp.StatusCode, string(body))
		return User{}, fmt.Errorf("failed to insert user: %s", body)
	}

	var insertedUsers []User
	err = json.Unmarshal(body, &insertedUsers)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return User{}, err
	}

	if len(insertedUsers) == 0 {
		return User{}, fmt.Errorf("no user was inserted")
	}

	insertedUser := insertedUsers[0]
	log.Printf("User inserted successfully: %s", insertedUser.Username)
	return insertedUser, nil
}

func GetAllUsers() ([]User, error) {
	url := fmt.Sprintf("%s/rest/v1/users", SupabaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", SupabaseKey)
	req.Header.Set("Authorization", "Bearer "+SupabaseKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get users: %s", body)
	}

	var users []User
	err = json.Unmarshal(body, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}
