package tests

import (
	"os"
	"testing"

	"github.com/revenium/revenium-middleware-runway-go/revenium"
	"github.com/stretchr/testify/assert"
)

func TestConfigLoadFromEnv(t *testing.T) {
	// Set test environment variables
	os.Setenv("RUNWAY_API_KEY", "test-runway-key")
	os.Setenv("REVENIUM_METERING_API_KEY", "hak_test-revenium-key")
	os.Setenv("RUNWAY_BASE_URL", "https://test.api.runwayml.com")
	os.Setenv("REVENIUM_LOG_LEVEL", "DEBUG")

	defer func() {
		os.Unsetenv("RUNWAY_API_KEY")
		os.Unsetenv("REVENIUM_METERING_API_KEY")
		os.Unsetenv("RUNWAY_BASE_URL")
		os.Unsetenv("REVENIUM_LOG_LEVEL")
	}()

	cfg := &revenium.Config{}
	err := cfg.LoadFromEnv()

	assert.NoError(t, err)
	assert.Equal(t, "test-runway-key", cfg.RunwayAPIKey)
	assert.Equal(t, "hak_test-revenium-key", cfg.ReveniumAPIKey)
	assert.Equal(t, "https://test.api.runwayml.com", cfg.RunwayBaseURL)
	assert.Equal(t, "DEBUG", cfg.LogLevel)
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      *revenium.Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			config: &revenium.Config{
				RunwayAPIKey:   "test-runway-key",
				ReveniumAPIKey: "hak_test-key",
			},
			expectError: false,
		},
		{
			name: "missing revenium key",
			config: &revenium.Config{
				RunwayAPIKey: "test-runway-key",
			},
			expectError: true,
			errorMsg:    "REVENIUM_METERING_API_KEY is required",
		},
		{
			name: "missing runway key",
			config: &revenium.Config{
				ReveniumAPIKey: "hak_test-key",
			},
			expectError: true,
			errorMsg:    "RUNWAY_API_KEY is required",
		},
		{
			name: "invalid revenium key format",
			config: &revenium.Config{
				RunwayAPIKey:   "test-runway-key",
				ReveniumAPIKey: "invalid-key-format",
			},
			expectError: true,
			errorMsg:    "invalid Revenium API key format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNormalizeReveniumBaseURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://api.revenium.ai", "https://api.revenium.ai"},
		{"https://api.revenium.ai/", "https://api.revenium.ai"},
		{"https://api.revenium.ai/meter/v2", "https://api.revenium.ai"},
		{"https://api.revenium.ai/meter", "https://api.revenium.ai"},
		{"https://api.revenium.ai/v2", "https://api.revenium.ai"},
		{"", "https://api.revenium.ai"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := revenium.NormalizeReveniumBaseURL(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfigOptions(t *testing.T) {
	cfg := &revenium.Config{}

	// Test WithRunwayAPIKey
	opt := revenium.WithRunwayAPIKey("test-key")
	opt(cfg)
	assert.Equal(t, "test-key", cfg.RunwayAPIKey)

	// Test WithReveniumAPIKey
	opt = revenium.WithReveniumAPIKey("hak_test")
	opt(cfg)
	assert.Equal(t, "hak_test", cfg.ReveniumAPIKey)

	// Test WithRunwayBaseURL
	opt = revenium.WithRunwayBaseURL("https://custom.url")
	opt(cfg)
	assert.Equal(t, "https://custom.url", cfg.RunwayBaseURL)

	// Test WithReveniumBaseURL
	opt = revenium.WithReveniumBaseURL("https://custom.revenium.url")
	opt(cfg)
	assert.Equal(t, "https://custom.revenium.url", cfg.ReveniumBaseURL)
}
