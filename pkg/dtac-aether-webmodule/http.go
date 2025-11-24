package dtac_aether_webmodule

import (
	"encoding/json"
	"net/http"
)

// RegisterRoutes registers the HTTP routes for this web module on the given mux.
//
// Adjust the paths or add more routes as needed for your module.
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", HealthHandler)
	mux.HandleFunc("/hello", HelloHandler)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	resp := map[string]string{
		"status": "ok",
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

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
