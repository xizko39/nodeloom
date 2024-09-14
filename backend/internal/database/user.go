package database

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID        string    `json:"id,omitempty"`
	Username  string    `json:"username"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// InsertUser inserts a new user into the Supabase database
// InsertUser inserts a new user into the Supabase database
func InsertUser(client *SupabaseClient, user User) (User, error) {
	url := "users" // Just the endpoint, not the full URL

	userData := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: user.Username,
		Password: user.Password,
	}

	body, status, err := client.Request("POST", url, userData) // Pass the client
	if err != nil {
		log.Printf("Error making request to Supabase: %v", err)
		return User{}, err
	}

	if status != http.StatusCreated {
		log.Printf("Supabase returned non-201 status: %d, body: %s", status, string(body))
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

	return insertedUsers[0], nil
}

// GetAllUsers retrieves all users from the Supabase database
func GetAllUsers() ([]User, error) {
	url := fmt.Sprintf("%s/rest/v1/users", SupabaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}

	req.Header.Set("apikey", SupabaseKey)
	req.Header.Set("Authorization", "Bearer "+SupabaseKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request to Supabase: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Supabase returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("failed to retrieve users: %s", body)
	}

	var users []User
	err = json.Unmarshal(body, &users)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}

	return users, nil
}

// FindUserByUsernameOrEmail searches for users by username or email in the Supabase database
func FindUserByUsernameOrEmail(client *SupabaseClient, identifier string) ([]User, error) {
	endpoint := "users?or=(username.eq." + identifier + ",email.eq." + identifier + ")"

	respBody, statusCode, err := client.Request("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error making request to Supabase: %v", err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("supabase returned non-200 status: %d, body: %s", statusCode, string(respBody))
	}

	var users []User
	err = json.Unmarshal(respBody, &users)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return users, nil
}
