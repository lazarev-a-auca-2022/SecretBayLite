package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"secretbay/internal/config"
)

func TestHealthCheck(t *testing.T) {
	cfg := &config.Config{
		ServerAddress: ":8080",
		JWTSecret:     "test-secret",
	}
	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)
	server := NewServer(cfg, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Health check failed: got %v want %v", w.Code, http.StatusOK)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if status, ok := response["status"]; !ok || status != "healthy" {
		t.Errorf("Unexpected health status: got %v want %v", status, "healthy")
	}
}

func TestVPNConfiguration(t *testing.T) {
	cfg := &config.Config{
		ServerAddress: ":8080",
		JWTSecret:     "test-secret",
	}
	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)
	server := NewServer(cfg, logger)

	// Test invalid request
	invalidReq := map[string]interface{}{
		"server_ip": "invalid-ip",
		"username":  "test",
	}
	body, _ := json.Marshal(invalidReq)

	req := httptest.NewRequest(http.MethodPost, "/api/vpn/configure", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected unauthorized without token: got %v want %v", w.Code, http.StatusUnauthorized)
	}

	// Test with invalid token
	req.Header.Set("Authorization", "Bearer invalid-token")
	w = httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected unauthorized with invalid token: got %v want %v", w.Code, http.StatusUnauthorized)
	}
}
