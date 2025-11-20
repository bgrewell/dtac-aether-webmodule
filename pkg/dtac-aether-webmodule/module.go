package {{module}}

import (
	"fmt"
	"net/http"
)

// ModuleName is the canonical name of this module.
const ModuleName = "{{module}}"

// HealthHandler is a simple health check endpoint.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// HelloHandler is a basic example handler for your web module.
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "world"
	}
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "Hello, %s from %s module!", name, ModuleName)
}
