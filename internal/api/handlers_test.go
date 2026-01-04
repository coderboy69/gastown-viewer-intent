package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/intent-solutions-io/gastown-viewer-intent/internal/beads"
)

func TestHealthHandler(t *testing.T) {
	// Create server with mock adapter
	config := DefaultConfig()
	config.TownRoot = "/tmp/nonexistent"
	adapter := beads.NewCLIAdapter("")

	server := NewServer(config, adapter)

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()

	server.Handler().ServeHTTP(w, req)

	// Health returns 200 if beads available, 503 if not (CI has no bd CLI)
	if w.Code != http.StatusOK && w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status 200 or 503, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Version should always be present
	if resp["version"] != config.Version {
		t.Errorf("Expected version %s, got %v", config.Version, resp["version"])
	}
}

func TestTownStatusHandler(t *testing.T) {
	config := DefaultConfig()
	config.TownRoot = "/tmp/nonexistent-town"
	adapter := beads.NewCLIAdapter("")

	server := NewServer(config, adapter)

	req := httptest.NewRequest("GET", "/api/v1/town/status", nil)
	w := httptest.NewRecorder()

	server.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Should report unhealthy for non-existent town
	if resp["healthy"] != false {
		t.Error("Expected healthy=false for non-existent town")
	}
}

func TestCORSMiddleware(t *testing.T) {
	config := DefaultConfig()
	config.CORSOrigins = []string{"http://localhost:5173"}
	adapter := beads.NewCLIAdapter("")

	server := NewServer(config, adapter)

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()

	server.Handler().ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Error("Expected CORS header to be set")
	}
}

func TestCORSPreflight(t *testing.T) {
	config := DefaultConfig()
	config.CORSOrigins = []string{"http://localhost:5173"}
	adapter := beads.NewCLIAdapter("")

	server := NewServer(config, adapter)

	req := httptest.NewRequest("OPTIONS", "/api/v1/health", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()

	server.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204 for preflight, got %d", w.Code)
	}
}
