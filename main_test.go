package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthRoute(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code, "HTTP status code should be 200")

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON response: %s", err)
	}
	assert.Equal(t, "healthy", response["status"], "Health status should be 'healthy'")
}
