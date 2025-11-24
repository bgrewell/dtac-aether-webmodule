# dtac-aether-webmodule

A modular web service component for the DTAC ecosystem that can run standalone or be embedded in dtac-agent.

## Overview

The dtac-aether-webmodule provides HTTP endpoints for the DTAC Aether system. It follows a clean module pattern with lifecycle management, making it suitable for both standalone execution and integration with the dtac-agent runtime.

## Features

- **Lifecycle Management**: Start/Stop methods with graceful shutdown
- **Environment-based Configuration**: Configure via environment variables
- **Health Checks**: Built-in `/healthz` and `/readyz` endpoints
- **Standalone or Embedded**: Run as a standalone service or embed in dtac-agent
- **Comprehensive Testing**: Full unit test coverage

## Running Standalone

### Quick Start

```bash
# Build the module
make build

# Run with defaults
./bin/dtac-aether-webmodule

# Or run directly with go
go run ./cmd/dtac-aether-webmodule
```

### Configuration

Configure the module using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `MODULE_LISTEN_ADDR` | Listen address and port | `:8080` |
| `MODULE_PREFIX` | URL prefix for module routes | `/aether` |
| `MODULE_READ_TIMEOUT` | HTTP read timeout (e.g., "30s") | `15s` |
| `MODULE_WRITE_TIMEOUT` | HTTP write timeout (e.g., "30s") | `15s` |

### Example

```bash
# Run on custom port with custom prefix
export MODULE_LISTEN_ADDR=":9090"
export MODULE_PREFIX="/api/v1"
export MODULE_READ_TIMEOUT="30s"
export MODULE_WRITE_TIMEOUT="30s"
./bin/dtac-aether-webmodule
```

## API Endpoints

### Health Endpoints

- **GET /healthz** - Health check endpoint
  ```bash
  curl http://localhost:8080/healthz
  # Response: {"status":"ok"}
  ```

- **GET /readyz** - Readiness check endpoint
  ```bash
  curl http://localhost:8080/readyz
  # Response: {"status":"ready"}
  ```

### Module Endpoints

All module-specific endpoints are prefixed with the configured `MODULE_PREFIX` (default: `/aether`).

- **GET /aether/** - Module information
  ```bash
  curl http://localhost:8080/aether/
  # Response: {"module":"dtac-aether-webmodule","version":"1.0.0","status":"running"}
  ```

- **GET /aether/hello** - Hello endpoint
  ```bash
  curl http://localhost:8080/aether/hello
  # Response: {"message":"hello from dtac-aether-webmodule"}
  ```

## Using as a Library

### Basic Usage

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/bgrewell/dtac-web-module-template/pkg/dtacaether"
)

func main() {
    // Load configuration from environment
    cfg, err := dtacaether.LoadConfigFromEnv()
    if err != nil {
        log.Fatal(err)
    }
    
    // Create module
    module := dtacaether.New(cfg)
    
    // Start the module
    ctx := context.Background()
    if err := module.Start(ctx); err != nil {
        log.Fatal(err)
    }
    
    // ... do other work ...
    
    // Stop the module gracefully
    stopCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := module.Stop(stopCtx); err != nil {
        log.Printf("error during shutdown: %v", err)
    }
}
```

### Using with Custom Configuration

```go
cfg := &dtacaether.Config{
    ListenAddr:   ":9090",
    Prefix:       "/myservice",
    ReadTimeout:  30 * time.Second,
    WriteTimeout: 30 * time.Second,
}

module := dtacaether.New(cfg)
```

### Integrating Routes with Existing ServeMux

```go
import (
    "net/http"
    "github.com/bgrewell/dtac-web-module-template/pkg/dtacaether"
)

func main() {
    // Create your own mux
    mux := http.NewServeMux()
    
    // Add your own routes
    mux.HandleFunc("/custom", myHandler)
    
    // Register dtac-aether routes (backwards compatibility)
    dtacaether.RegisterRoutes(mux)
    
    // Or create a module and register its routes
    cfg := &dtacaether.Config{Prefix: "/aether", ListenAddr: ":8080"}
    module := dtacaether.New(cfg)
    module.RegisterRoutes(mux)
    
    // Use the mux with your server
    http.ListenAndServe(":8080", mux)
}
```

## Embedding in dtac-agent

The module can be registered with dtac-agent for managed execution. The existing `pkg/dtac-aether-webmodule` package provides the dtac-agent Module interface implementation.

```go
import (
    "github.com/bgrewell/dtac-web-module-template/pkg/dtac-aether-webmodule"
)

// Register with dtac-agent
module := dtac_aether_webmodule.NewAetherWebModule()
// ... register with dtac-agent module manager
```

## Development

### Building

```bash
# Build the binary
make build

# Clean build artifacts
make clean
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./pkg/dtacaether/...

# Run tests with coverage
go test -cover ./pkg/dtacaether/...
```

### Project Structure

```
.
├── cmd/
│   └── dtac-aether-webmodule/
│       ├── main.go           # Standalone bootstrap
│       └── build.yaml        # Build configuration
├── pkg/
│   ├── dtacaether/          # Standalone module package
│   │   ├── module.go        # Module lifecycle implementation
│   │   ├── config.go        # Configuration management
│   │   ├── handlers.go      # HTTP handlers
│   │   ├── health.go        # Health check handlers
│   │   └── *_test.go        # Unit tests
│   └── dtac-aether-webmodule/ # dtac-agent integration
│       ├── module.go        # Module interface implementation
│       └── http.go          # Original HTTP handlers
├── Makefile                 # Build automation
└── README.md               # This file
```

## Graceful Shutdown

The module supports graceful shutdown via context cancellation. When running standalone, the bootstrap handles SIGINT and SIGTERM signals:

```bash
# Start the module
./bin/dtac-aether-webmodule

# Press Ctrl+C to trigger graceful shutdown
# The module will:
# 1. Stop accepting new connections
# 2. Wait for active requests to complete (up to 30s)
# 3. Shut down cleanly
```

## License

See [LICENSE](../../LICENSE) for details.

## Contributing

Contributions are welcome! Please ensure:
1. All tests pass: `go test ./...`
2. Code is formatted: `go fmt ./...`
3. Code is linted: `go vet ./...`
