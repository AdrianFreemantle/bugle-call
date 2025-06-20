// Package config handles configuration management for the async processor service.
// It loads and validates configuration from environment variables.
package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

// Variable for testing - allows us to mock os.Exit
var osExit = os.Exit

// Config holds all configuration for the async processor service.
type Config struct {
	// HTTP server configuration
	HTTPPort string `envconfig:"HTTP_PORT" default:"8080"`

	// Logging configuration
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`

	// NATS configuration
	NATSURL string `envconfig:"NATS_URL" required:"true"`

	// MongoDB configuration (placeholder for future use)
	MongoURI string `envconfig:"MONGO_URI" default:""`
}

// New creates a new Config instance by loading values from environment variables.
// Returns an error if required environment variables are missing or invalid.
func New() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Validate NATS URL
	if cfg.NATSURL == "" {
		return nil, fmt.Errorf("NATS_URL cannot be empty")
	}
	
	if _, err := url.ParseRequestURI(cfg.NATSURL); err != nil {
		return nil, fmt.Errorf("invalid NATS_URL: %w", err)
	}

	// Normalize and validate log level
	cfg.LogLevel = strings.ToLower(cfg.LogLevel)
	switch cfg.LogLevel {
	case "debug", "info", "warn", "error":
		// Valid log level, keep as is
	default:
		// Invalid log level, fallback to info
		cfg.LogLevel = "info"
	}

	return &cfg, nil
}

// MustLoad loads the configuration and exits the program if any error occurs.
// This is intended for use during application startup.
func MustLoad() *Config {
	cfg, err := New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: %v\n", err)
		osExit(1)
	}
	return cfg
}
