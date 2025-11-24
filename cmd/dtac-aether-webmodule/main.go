package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bgrewell/dtac-web-module-template/pkg/dtacaether"
)

const (
	// shutdownTimeout is the maximum time allowed for graceful shutdown.
	shutdownTimeout = 30 * time.Second
)

func main() {
	// Load configuration from environment variables
	cfg, err := dtacaether.LoadConfigFromEnv()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Create the module
	module := dtacaether.New(cfg)

	// Create a context that will be canceled on interrupt signal
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start the module in a goroutine
	if err := module.Start(ctx); err != nil {
		log.Fatalf("failed to start module: %v", err)
	}

	// Wait for interrupt signal
	<-sigChan
	log.Println("Received shutdown signal, stopping gracefully...")

	// Create a shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	// Stop the module
	if err := module.Stop(shutdownCtx); err != nil {
		log.Printf("error during shutdown: %v", err)
		os.Exit(1)
	}

	log.Println("Shutdown complete")
}
