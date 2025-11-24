package dtacaether

import (
	"encoding/json"
	"net/http"
)

// HelloHandler handles the /hello endpoint.
// It returns a JSON response with a greeting message.
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	resp := map[string]string{
		"message": "hello from dtac-aether-webmodule",
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// RootHandler handles the root endpoint.
// It returns a JSON response with module information.
func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	resp := map[string]interface{}{
		"module":  "dtac-aether-webmodule",
		"version": "1.0.0",
		"status":  "running",
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
