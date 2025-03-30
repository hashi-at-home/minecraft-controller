package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test for health endpoint
func TestHealth(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// Test for Readiness Endpoint
// No tokens should return 200 with ready: false
func TestReadiness(t *testing.T) {
	router := setupRouter()
	router = setupReadiness(router)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/readiness", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var r Ready
	err := json.NewDecoder(w.Body).Decode(&r)
	assert.NoError(t, err)
	assert.False(t, r.Ready)
	assert.False(t, r.DigitalOceanInitialized)
}

func TestRoot(t *testing.T) {
	router := setupRouter()
	router = setupStaticAssets(router)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
