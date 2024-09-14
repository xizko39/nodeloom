// internal/database/supabase_client.go

package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// SupabaseClient encapsulates the Supabase configuration and the HTTP client
type SupabaseClient struct {
	URL  string
	Key  string
	HTTP *http.Client
}

// NewSupabaseClient initializes a new Supabase client
func NewSupabaseClient(url, key string) *SupabaseClient {
	if url == "" || key == "" {
		fmt.Printf("Supabase URL or Key is missing! URL: %s, Key: %s\n", url, key)
		return nil
	}
	return &SupabaseClient{
		URL: url,
		Key: key,
		HTTP: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Request performs a generic HTTP request to Supabase
func (c *SupabaseClient) Request(method, endpoint string, body interface{}) ([]byte, int, error) {
	if c.URL == "" {
		return nil, 0, fmt.Errorf("supabase URL is empty")
	}

	fullURL := fmt.Sprintf("%s/rest/v1/%s", c.URL, endpoint)
	fmt.Printf("Full URL: %s\n", fullURL)

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return nil, 0, err
	}

	// Set headers
	req.Header.Set("apikey", c.Key)
	req.Header.Set("Authorization", "Bearer "+c.Key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return respBody, resp.StatusCode, nil
}
