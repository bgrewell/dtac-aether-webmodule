package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bgrewell/dtac-web-module-template/pkg/dtac-aether-webmodule"
)

func main() {
	// Example: use an environment variable or config for listen address
	addr := os.Getenv("MODULE_LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	mux := http.NewServeMux()

	// Register HTTP routes for this module.
	dtac_aether_webmodule.RegisterRoutes(mux)

	log.Printf("Starting dtac-aether-webmodule web module on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("failed to start dtac-aether-webmodule web module: %v", err)
	}
}
