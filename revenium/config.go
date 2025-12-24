package revenium

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the Revenium middleware
type Config struct {
	// Runway API configuration
	RunwayAPIKey string
	RunwayBaseURL string
	RunwayVersion string

	// Revenium metering configuration
	ReveniumAPIKey    string
	ReveniumBaseURL   string
	ReveniumOrgID     string
	ReveniumProductID string

	// Logging and debug configuration
	LogLevel       string
	VerboseStartup bool
}

// Option is a functional option for configuring Config
type Option func(*Config)

// WithRunwayAPIKey sets the Runway API key
func WithRunwayAPIKey(key string) Option {
	return func(c *Config) {
		c.RunwayAPIKey = key
	}
}

// WithRunwayBaseURL sets the Runway base URL
func WithRunwayBaseURL(url string) Option {
	return func(c *Config) {
		c.RunwayBaseURL = url
	}
}

// WithReveniumAPIKey sets the Revenium API key
func WithReveniumAPIKey(key string) Option {
	return func(c *Config) {
		c.ReveniumAPIKey = key
	}
}

// WithReveniumBaseURL sets the Revenium base URL
func WithReveniumBaseURL(url string) Option {
	return func(c *Config) {
		c.ReveniumBaseURL = url
	}
}

// LoadFromEnv loads configuration from environment variables and .env files
func (c *Config) LoadFromEnv() error {
	// First, try to load .env files automatically
	c.loadEnvFiles()

	// Then load from environment variables (which may have been set by .env files)
	c.RunwayAPIKey = os.Getenv("RUNWAY_API_KEY")
	c.RunwayBaseURL = getEnvOrDefault("RUNWAY_BASE_URL", "https://api.runwayml.com")
	c.RunwayVersion = getEnvOrDefault("RUNWAY_VERSION", "2024-11-06")

	c.ReveniumAPIKey = os.Getenv("REVENIUM_METERING_API_KEY")
	baseURL := getEnvOrDefault("REVENIUM_METERING_BASE_URL", "https://api.revenium.ai")
	c.ReveniumBaseURL = NormalizeReveniumBaseURL(baseURL)
	c.ReveniumOrgID = os.Getenv("REVENIUM_ORGANIZATION_ID")
	c.ReveniumProductID = os.Getenv("REVENIUM_PRODUCT_ID")

	c.LogLevel = getEnvOrDefault("REVENIUM_LOG_LEVEL", "INFO")
	c.VerboseStartup = os.Getenv("REVENIUM_VERBOSE_STARTUP") == "true" || os.Getenv("REVENIUM_VERBOSE_STARTUP") == "1"

	// Initialize logger early so we can use it
	InitializeLogger()

	// Debug log for configuration loading
	Debug("Loading configuration from environment variables")
	if c.RunwayAPIKey != "" {
		Debug("Runway API key loaded (length: %d)", len(c.RunwayAPIKey))
	}

	return nil
}

// loadEnvFiles loads environment variables from .env files
func (c *Config) loadEnvFiles() {
	// Try to load .env files in order of preference
	envFiles := []string{
		".env.local", // Local overrides (highest priority)
		".env",       // Main env file
	}

	var loadedFiles []string

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	// Try current directory and parent directories
	searchDirs := []string{
		cwd,
		filepath.Dir(cwd),
		filepath.Join(cwd, ".."),
	}

	for _, dir := range searchDirs {
		for _, envFile := range envFiles {
			envPath := filepath.Join(dir, envFile)

			// Check if file exists
			if _, err := os.Stat(envPath); err == nil {
				// Try to load the file
				if err := godotenv.Load(envPath); err == nil {
					loadedFiles = append(loadedFiles, envPath)
				}
			}
		}
	}

	// Log loaded files (only if we have a logger initialized)
	if len(loadedFiles) > 0 {
		// We can't use Debug here because logger might not be initialized yet
		// So we'll just silently load the files
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.ReveniumAPIKey == "" {
		return NewConfigError("REVENIUM_METERING_API_KEY is required", nil)
	}

	if !isValidAPIKeyFormat(c.ReveniumAPIKey) {
		return NewConfigError("invalid Revenium API key format", nil)
	}

	if c.RunwayAPIKey == "" {
		return NewConfigError("RUNWAY_API_KEY is required", nil)
	}

	Debug("Configuration validation passed")
	return nil
}

// isValidAPIKeyFormat checks if the API key has a valid format
func isValidAPIKeyFormat(key string) bool {
	// Revenium API keys should start with "hak_"
	if len(key) < 4 {
		return false
	}
	return key[:4] == "hak_"
}

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// NormalizeReveniumBaseURL normalizes the base URL to a consistent format
// It handles various input formats and returns a normalized base URL without trailing slash
func NormalizeReveniumBaseURL(baseURL string) string {
	if baseURL == "" {
		return "https://api.revenium.ai"
	}

	// Remove trailing slash if present
	if len(baseURL) > 0 && baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}

	// If it already ends with /meter/v2, remove /meter/v2 (legacy format)
	if len(baseURL) >= 9 && baseURL[len(baseURL)-9:] == "/meter/v2" {
		return baseURL[:len(baseURL)-9]
	}

	// If it ends with /meter, remove /meter (legacy format)
	if len(baseURL) >= 6 && baseURL[len(baseURL)-6:] == "/meter" {
		return baseURL[:len(baseURL)-6]
	}

	// If it ends with /v2, remove /v2 (legacy format)
	if len(baseURL) >= 3 && baseURL[len(baseURL)-3:] == "/v2" {
		return baseURL[:len(baseURL)-3]
	}

	// Return the base URL as-is (should be just the domain)
	return baseURL
}
