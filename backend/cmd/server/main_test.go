package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/xizko39/nodeloom/internal/api/routes"
	"github.com/xizko39/nodeloom/internal/config"
	"github.com/xizko39/nodeloom/internal/database"
)

func setupRouter() *gin.Engine {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	database.InitSupabase(cfg.Supabase.URL, cfg.Supabase.Key)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	routes.SetupRoutes(r)
	return r
}

func TestRegisterUser(t *testing.T) {
	router := setupRouter()

	// Create a test user
	user := map[string]string{
		"username": "testuser",
		"password": "testpassword",
	}
	jsonValue, _ := json.Marshal(user)

	// Create a new request
	req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check the status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d; got %d", http.StatusCreated, w.Code)
	}

	// Parse the response body
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	// Check the response message
	expectedMessage := "User registered successfully"
	if message, exists := response["message"]; !exists || message != expectedMessage {
		t.Errorf("Expected message '%s'; got '%s'", expectedMessage, message)
	}

	// Check if the username is returned in the response
	if username, exists := response["username"]; !exists || username != user["username"] {
		t.Errorf("Expected username '%s'; got '%s'", user["username"], username)
	}
}
