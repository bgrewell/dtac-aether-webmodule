package dtacaether

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	// startupCheckDelay is the time to wait after starting the server to check for immediate failures.
	startupCheckDelay = 100 * time.Millisecond
)

// Module represents the dtac-aether-webmodule with lifecycle management.
type Module struct {
	cfg    *Config
	server *http.Server
	mux    *http.ServeMux
	logger *log.Logger
	mu     sync.Mutex
	
	// shutdown is closed when the module is stopped
	shutdown chan struct{}
	// serverErr receives any error from the server's ListenAndServe call
	serverErr chan error
}

// New creates a new Module instance with the given configuration.
// If cfg is nil, it uses a default configuration.
func New(cfg *Config) *Module {
	if cfg == nil {
		cfg = &Config{
			ListenAddr:   ":8080",
			Prefix:       "/aether",
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		}
	}
	
	mux := http.NewServeMux()
	
	m := &Module{
		cfg:       cfg,
		mux:       mux,
		logger:    log.Default(),
		shutdown:  make(chan struct{}),
		serverErr: make(chan error, 1),
	}
	
	// Register routes
	m.RegisterRoutes(mux)
	
	// Create the HTTP server
	m.server = &http.Server{
		Addr:         cfg.ListenAddr,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
	
	return m
}

// Name returns the name of the module.
func (m *Module) Name() string {
	return "dtac-aether-webmodule"
}

// RegisterRoutes registers the HTTP routes for this module.
// Routes are prefixed with the configured prefix.
func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	prefix := m.cfg.Prefix
	
	// Ensure prefix starts with / and doesn't end with /
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	prefix = strings.TrimSuffix(prefix, "/")
	
	// Register health endpoints (no prefix for standard health checks)
	mux.HandleFunc("/healthz", HealthHandler)
	mux.HandleFunc("/readyz", ReadyHandler)
	
	// Register module endpoints with prefix
	mux.HandleFunc(prefix+"/", RootHandler)
	mux.HandleFunc(prefix+"/hello", HelloHandler)
}

// Routes returns the underlying http.Handler (ServeMux) for the module.
// This can be used to integrate the module with other HTTP routers.
func (m *Module) Routes() http.Handler {
	return m.mux
}

// Start starts the HTTP server in a goroutine.
// It returns immediately after launching the server goroutine.
// The context can be used to cancel the start operation if needed.
// Returns an error if the server fails to start.
func (m *Module) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Check if already started
	select {
	case <-m.shutdown:
		return fmt.Errorf("module has been shut down")
	default:
	}
	
	m.logger.Printf("Starting %s on %s", m.Name(), m.cfg.ListenAddr)
	
	// Start server in goroutine
	go func() {
		if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			m.logger.Printf("Server error: %v", err)
			m.serverErr <- err
		}
	}()
	
	// Give the server a moment to start and potentially fail
	select {
	case err := <-m.serverErr:
		return fmt.Errorf("failed to start server: %w", err)
	case <-time.After(startupCheckDelay):
		m.logger.Printf("%s started successfully on %s", m.Name(), m.cfg.ListenAddr)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Stop gracefully shuts down the HTTP server.
// It waits for active connections to complete or for the context to be canceled.
func (m *Module) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Check if already shut down
	select {
	case <-m.shutdown:
		return nil
	default:
		close(m.shutdown)
	}
	
	m.logger.Printf("Stopping %s...", m.Name())
	
	if m.server == nil {
		return nil
	}
	
	if err := m.server.Shutdown(ctx); err != nil {
		m.logger.Printf("Error during shutdown: %v", err)
		return fmt.Errorf("failed to shutdown server: %w", err)
	}
	
	m.logger.Printf("%s stopped successfully", m.Name())
	return nil
}

// SetLogger sets a custom logger for the module.
func (m *Module) SetLogger(logger *log.Logger) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logger = logger
}

// RegisterRoutes is a package-level convenience function that creates a default module
// and registers its routes on the provided mux.
// This provides backwards compatibility with code that uses the old RegisterRoutes pattern.
func RegisterRoutes(mux *http.ServeMux) {
	cfg := &Config{
		ListenAddr:   ":8080",
		Prefix:       "/aether",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	
	m := New(cfg)
	m.RegisterRoutes(mux)
}
