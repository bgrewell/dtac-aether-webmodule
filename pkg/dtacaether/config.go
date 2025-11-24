package dtacaether

import (
	"os"
	"strconv"
	"time"
)

// Config holds the configuration for the Aether web module.
type Config struct {
	// ListenAddr is the address to listen on (e.g., ":8080").
	ListenAddr string
	
	// Prefix is the URL prefix for all routes (e.g., "/aether").
	Prefix string
	
	// ReadTimeout is the maximum duration for reading the entire request.
	ReadTimeout time.Duration
	
	// WriteTimeout is the maximum duration before timing out writes of the response.
	WriteTimeout time.Duration
}

// LoadConfigFromEnv loads configuration from environment variables with sensible defaults.
// Environment variables:
//   - MODULE_LISTEN_ADDR: Listen address (default: ":8080")
//   - MODULE_PREFIX: URL prefix (default: "/aether")
//   - MODULE_READ_TIMEOUT: Read timeout duration (default: "15s")
//   - MODULE_WRITE_TIMEOUT: Write timeout duration (default: "15s")
func LoadConfigFromEnv() (*Config, error) {
	cfg := &Config{
		ListenAddr:   getEnvOrDefault("MODULE_LISTEN_ADDR", ":8080"),
		Prefix:       getEnvOrDefault("MODULE_PREFIX", "/aether"),
		ReadTimeout:  parseDurationOrDefault("MODULE_READ_TIMEOUT", 15*time.Second),
		WriteTimeout: parseDurationOrDefault("MODULE_WRITE_TIMEOUT", 15*time.Second),
	}
	
	return cfg, nil
}

// getEnvOrDefault returns the value of an environment variable or a default value.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseDurationOrDefault parses a duration from an environment variable or returns a default.
func parseDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	
	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	
	return duration
}

// parseIntOrDefault parses an integer from an environment variable or returns a default.
func parseIntOrDefault(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	
	return intVal
}
