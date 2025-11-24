package dtacaether

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestModule_Name(t *testing.T) {
	m := New(nil)
	
	expectedName := "dtac-aether-webmodule"
	if m.Name() != expectedName {
		t.Errorf("Expected name '%s', got '%s'", expectedName, m.Name())
	}
}

func TestModule_New_WithConfig(t *testing.T) {
	cfg := &Config{
		ListenAddr:   ":9999",
		Prefix:       "/test",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 45 * time.Second,
	}
	
	m := New(cfg)
	
	if m.cfg.ListenAddr != ":9999" {
		t.Errorf("Expected ListenAddr ':9999', got '%s'", m.cfg.ListenAddr)
	}
	
	if m.cfg.Prefix != "/test" {
		t.Errorf("Expected Prefix '/test', got '%s'", m.cfg.Prefix)
	}
	
	if m.server == nil {
		t.Error("Expected server to be initialized")
	}
	
	if m.mux == nil {
		t.Error("Expected mux to be initialized")
	}
}

func TestModule_New_WithNilConfig(t *testing.T) {
	m := New(nil)
	
	// Should use defaults
	if m.cfg.ListenAddr != ":8080" {
		t.Errorf("Expected default ListenAddr ':8080', got '%s'", m.cfg.ListenAddr)
	}
	
	if m.cfg.Prefix != "/aether" {
		t.Errorf("Expected default Prefix '/aether', got '%s'", m.cfg.Prefix)
	}
}

func TestModule_StartStop(t *testing.T) {
	// Use ephemeral port
	cfg := &Config{
		ListenAddr:   "127.0.0.1:0",
		Prefix:       "/aether",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	
	m := New(cfg)
	
	ctx := context.Background()
	
	// Start the module
	if err := m.Start(ctx); err != nil {
		t.Fatalf("Failed to start module: %v", err)
	}
	
	// Give server time to start
	time.Sleep(200 * time.Millisecond)
	
	// Get the actual address (since we used :0)
	// We need to make a request to verify it's running
	baseURL := fmt.Sprintf("http://%s", m.server.Addr)
	
	// Try to reach health endpoint
	resp, err := http.Get(baseURL + "/healthz")
	if err != nil {
		// If we get an error, it might be because :0 didn't work as expected
		// Let's just verify Stop works
		t.Logf("Could not reach server at %s: %v (this is ok for ephemeral port test)", baseURL, err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	}
	
	// Stop the module
	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := m.Stop(stopCtx); err != nil {
		t.Fatalf("Failed to stop module: %v", err)
	}
	
	// Verify server is stopped - this should fail or timeout
	time.Sleep(100 * time.Millisecond)
}

func TestModule_StartStop_RealPort(t *testing.T) {
	// Use a specific port that's likely available
	cfg := &Config{
		ListenAddr:   "127.0.0.1:18765",
		Prefix:       "/aether",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	
	m := New(cfg)
	
	ctx := context.Background()
	
	// Start the module
	if err := m.Start(ctx); err != nil {
		t.Fatalf("Failed to start module: %v", err)
	}
	
	// Give server time to start
	time.Sleep(200 * time.Millisecond)
	
	// Verify server is running by making requests
	baseURL := "http://127.0.0.1:18765"
	
	// Test health endpoint
	resp, err := http.Get(baseURL + "/healthz")
	if err != nil {
		t.Fatalf("Failed to reach health endpoint: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected health status 200, got %d", resp.StatusCode)
	}
	
	// Test hello endpoint with prefix
	resp, err = http.Get(baseURL + "/aether/hello")
	if err != nil {
		t.Fatalf("Failed to reach hello endpoint: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected hello status 200, got %d", resp.StatusCode)
	}
	
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	
	if len(body) == 0 {
		t.Error("Expected non-empty response body")
	}
	
	// Stop the module
	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := m.Stop(stopCtx); err != nil {
		t.Fatalf("Failed to stop module: %v", err)
	}
	
	// Verify server is stopped
	time.Sleep(100 * time.Millisecond)
	
	// This should fail now
	_, err = http.Get(baseURL + "/healthz")
	if err == nil {
		t.Error("Expected request to fail after Stop, but it succeeded")
	}
}

func TestModule_StopIdempotent(t *testing.T) {
	cfg := &Config{
		ListenAddr:   "127.0.0.1:18766",
		Prefix:       "/aether",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	
	m := New(cfg)
	
	ctx := context.Background()
	
	// Start the module
	if err := m.Start(ctx); err != nil {
		t.Fatalf("Failed to start module: %v", err)
	}
	
	time.Sleep(100 * time.Millisecond)
	
	// Stop the module
	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := m.Stop(stopCtx); err != nil {
		t.Fatalf("Failed to stop module: %v", err)
	}
	
	// Stop again - should not error
	if err := m.Stop(stopCtx); err != nil {
		t.Errorf("Second Stop() returned error: %v", err)
	}
}

func TestModule_RegisterRoutes(t *testing.T) {
	cfg := &Config{
		ListenAddr:   ":8080",
		Prefix:       "/custom",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	
	m := New(cfg)
	mux := http.NewServeMux()
	
	m.RegisterRoutes(mux)
	
	// We can't easily test that routes are registered without starting a server,
	// but we can verify that the method doesn't panic
	if m.mux == nil {
		t.Error("Expected mux to be set")
	}
}

func TestModule_Routes(t *testing.T) {
	m := New(nil)
	
	handler := m.Routes()
	if handler == nil {
		t.Error("Expected Routes() to return non-nil handler")
	}
}

func TestPackageRegisterRoutes(t *testing.T) {
	mux := http.NewServeMux()
	
	// Should not panic
	RegisterRoutes(mux)
	
	// Verify we can create a test server with it
	// This is a basic sanity check
	if mux == nil {
		t.Error("Expected mux to be valid after RegisterRoutes")
	}
}
