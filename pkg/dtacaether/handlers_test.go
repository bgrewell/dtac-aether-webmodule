package dtacaether

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name:           "GET request succeeds",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"status": "ok"},
		},
		{
			name:           "POST request fails",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   nil,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/healthz", nil)
			w := httptest.NewRecorder()
			
			HealthHandler(w, req)
			
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			
			if tt.expectedBody != nil {
				var resp map[string]string
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				
				if resp["status"] != tt.expectedBody["status"] {
					t.Errorf("Expected status '%s', got '%s'", tt.expectedBody["status"], resp["status"])
				}
			}
		})
	}
}

func TestReadyHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name:           "GET request succeeds",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"status": "ready"},
		},
		{
			name:           "PUT request fails",
			method:         http.MethodPut,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   nil,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/readyz", nil)
			w := httptest.NewRecorder()
			
			ReadyHandler(w, req)
			
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			
			if tt.expectedBody != nil {
				var resp map[string]string
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				
				if resp["status"] != tt.expectedBody["status"] {
					t.Errorf("Expected status '%s', got '%s'", tt.expectedBody["status"], resp["status"])
				}
			}
		})
	}
}

func TestHelloHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		checkMessage   bool
	}{
		{
			name:           "GET request succeeds",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			checkMessage:   true,
		},
		{
			name:           "POST request fails",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
			checkMessage:   false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/hello", nil)
			w := httptest.NewRecorder()
			
			HelloHandler(w, req)
			
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			
			if tt.checkMessage {
				var resp map[string]string
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				
				expectedMsg := "hello from dtac-aether-webmodule"
				if resp["message"] != expectedMsg {
					t.Errorf("Expected message '%s', got '%s'", expectedMsg, resp["message"])
				}
			}
		})
	}
}

func TestRootHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		checkBody      bool
	}{
		{
			name:           "GET request succeeds",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			checkBody:      true,
		},
		{
			name:           "DELETE request fails",
			method:         http.MethodDelete,
			expectedStatus: http.StatusMethodNotAllowed,
			checkBody:      false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", nil)
			w := httptest.NewRecorder()
			
			RootHandler(w, req)
			
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			
			if tt.checkBody {
				var resp map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				
				if resp["module"] != "dtac-aether-webmodule" {
					t.Errorf("Expected module 'dtac-aether-webmodule', got '%v'", resp["module"])
				}
				
				if resp["status"] != "running" {
					t.Errorf("Expected status 'running', got '%v'", resp["status"])
				}
			}
		})
	}
}
