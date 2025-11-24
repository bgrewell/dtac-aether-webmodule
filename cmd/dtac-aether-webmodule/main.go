package main

import (
	"github.com/bgrewell/dtac-agent/pkg/modules"
	"log"
	"os"

	"github.com/bgrewell/dtac-web-module-template/pkg/dtac-aether-webmodule"
)

func main() {

	// TODO: Need a way to feed the env vars for the web application through to it. A .env file is used with node but
	//       we will need to do something differently here I assume. Until this is done the web UI won't be able to find
	//       the backend server.

	// Example: use an environment variable or config for listen address
	addr := os.Getenv("MODULE_LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	awm := dtac_aether_webmodule.NewAetherWebModule()
	h, err := modules.NewModuleHost(awm)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting dtac-aether-webmodule web module on %s", addr)
	if err := h.Serve(); err != nil {
		log.Fatalf("failed to start dtac-aether-webmodule web module: %v", err)
	}
}
