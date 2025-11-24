package dtacaether

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfigFromEnv_Defaults(t *testing.T) {
	// Clear any existing env vars
	os.Unsetenv("MODULE_LISTEN_ADDR")
	os.Unsetenv("MODULE_PREFIX")
	os.Unsetenv("MODULE_READ_TIMEOUT")
	os.Unsetenv("MODULE_WRITE_TIMEOUT")
	
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("LoadConfigFromEnv() failed: %v", err)
	}
	
	if cfg.ListenAddr != ":8080" {
		t.Errorf("Expected ListenAddr to be ':8080', got '%s'", cfg.ListenAddr)
	}
	
	if cfg.Prefix != "/aether" {
		t.Errorf("Expected Prefix to be '/aether', got '%s'", cfg.Prefix)
	}
	
	if cfg.ReadTimeout != 15*time.Second {
		t.Errorf("Expected ReadTimeout to be 15s, got %v", cfg.ReadTimeout)
	}
	
	if cfg.WriteTimeout != 15*time.Second {
		t.Errorf("Expected WriteTimeout to be 15s, got %v", cfg.WriteTimeout)
	}
}

func TestLoadConfigFromEnv_CustomValues(t *testing.T) {
	// Set custom env vars
	os.Setenv("MODULE_LISTEN_ADDR", ":9090")
	os.Setenv("MODULE_PREFIX", "/custom")
	os.Setenv("MODULE_READ_TIMEOUT", "30s")
	os.Setenv("MODULE_WRITE_TIMEOUT", "45s")
	
	defer func() {
		os.Unsetenv("MODULE_LISTEN_ADDR")
		os.Unsetenv("MODULE_PREFIX")
		os.Unsetenv("MODULE_READ_TIMEOUT")
		os.Unsetenv("MODULE_WRITE_TIMEOUT")
	}()
	
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("LoadConfigFromEnv() failed: %v", err)
	}
	
	if cfg.ListenAddr != ":9090" {
		t.Errorf("Expected ListenAddr to be ':9090', got '%s'", cfg.ListenAddr)
	}
	
	if cfg.Prefix != "/custom" {
		t.Errorf("Expected Prefix to be '/custom', got '%s'", cfg.Prefix)
	}
	
	if cfg.ReadTimeout != 30*time.Second {
		t.Errorf("Expected ReadTimeout to be 30s, got %v", cfg.ReadTimeout)
	}
	
	if cfg.WriteTimeout != 45*time.Second {
		t.Errorf("Expected WriteTimeout to be 45s, got %v", cfg.WriteTimeout)
	}
}

func TestLoadConfigFromEnv_InvalidDuration(t *testing.T) {
	// Set invalid duration
	os.Setenv("MODULE_READ_TIMEOUT", "invalid")
	os.Setenv("MODULE_WRITE_TIMEOUT", "also-invalid")
	
	defer func() {
		os.Unsetenv("MODULE_READ_TIMEOUT")
		os.Unsetenv("MODULE_WRITE_TIMEOUT")
	}()
	
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("LoadConfigFromEnv() failed: %v", err)
	}
	
	// Should fall back to defaults on parse error
	if cfg.ReadTimeout != 15*time.Second {
		t.Errorf("Expected ReadTimeout to fall back to 15s, got %v", cfg.ReadTimeout)
	}
	
	if cfg.WriteTimeout != 15*time.Second {
		t.Errorf("Expected WriteTimeout to fall back to 15s, got %v", cfg.WriteTimeout)
	}
}

func TestParseDurationOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue time.Duration
		expected     time.Duration
	}{
		{
			name:         "valid duration",
			envKey:       "TEST_DURATION",
			envValue:     "5m",
			defaultValue: 1 * time.Second,
			expected:     5 * time.Minute,
		},
		{
			name:         "invalid duration uses default",
			envKey:       "TEST_DURATION",
			envValue:     "invalid",
			defaultValue: 10 * time.Second,
			expected:     10 * time.Second,
		},
		{
			name:         "empty value uses default",
			envKey:       "TEST_DURATION",
			envValue:     "",
			defaultValue: 20 * time.Second,
			expected:     20 * time.Second,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			}
			
			result := parseDurationOrDefault(tt.envKey, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue string
		expected     string
	}{
		{
			name:         "env var set",
			envKey:       "TEST_VAR",
			envValue:     "test-value",
			defaultValue: "default",
			expected:     "test-value",
		},
		{
			name:         "env var not set",
			envKey:       "TEST_VAR",
			envValue:     "",
			defaultValue: "default",
			expected:     "default",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			} else {
				os.Unsetenv(tt.envKey)
			}
			
			result := getEnvOrDefault(tt.envKey, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
