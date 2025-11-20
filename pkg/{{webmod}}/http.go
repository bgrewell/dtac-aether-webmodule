package {{module}}

import "net/http"

// RegisterRoutes registers the HTTP routes for this web module on the given mux.
//
// Adjust the paths or add more routes as needed for your module.
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", HealthHandler)
	mux.HandleFunc("/hello", HelloHandler)
}
