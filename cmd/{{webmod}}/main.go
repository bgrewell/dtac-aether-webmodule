package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bgrewell/dtac-web-module-template/pkg/{{module}}"
)

func main() {
	// Example: use an environment variable or config for listen address
	addr := os.Getenv("MODULE_LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	mux := http.NewServeMux()

	// Register HTTP routes for this module.
	{{module}}.RegisterRoutes(mux)

	log.Printf("Starting {{module}} web module on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("failed to start {{module}} web module: %v", err)
	}
}
