// Package config provides configuration management for the SecretBay application.
package config

import "os"

// Config holds the application configuration
type Config struct {
	// ServerAddress is the address and port the server listens on
	ServerAddress string `json:"server_address"`

	// JWTSecret is the secret key used for JWT token generation and validation
	JWTSecret string `json:"jwt_secret"`

	// LogPath specifies the directory where log files are stored
	LogPath string `json:"log_path"`
}

// Load reads the configuration from environment variables or sets defaults
func Load() (*Config, error) {
	cfg := &Config{
		ServerAddress: getEnvOrDefault("SERVER_ADDRESS", ":8080"),
		JWTSecret:     getEnvOrDefault("JWT_SECRET", "change_this_in_production"),
		LogPath:       getEnvOrDefault("LOG_PATH", "/app/logs"),
	}

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(cfg.LogPath, 0755); err != nil {
		return nil, err
	}

	return cfg, nil
}

// getEnvOrDefault returns the value of an environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
