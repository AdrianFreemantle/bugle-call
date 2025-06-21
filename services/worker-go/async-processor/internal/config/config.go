// Package config handles configuration management for the async processor service.
// It loads and validates configuration from environment variables.
package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

// ErrNATSURLRequired is returned when the NATS_URL is not provided.
var ErrNATSURLRequired = errors.New("NATS_URL is required and cannot be empty")

// Config holds all configuration for the async processor service.
type Config struct {
	// HTTPPort is the port the HTTP server listens on.
	HTTPPort string `envconfig:"HTTP_PORT" default:"8080"`

	// LogLevel is the logging level (e.g., "debug", "info", "warn", "error").
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`

	// NATSURL is the connection URL for the NATS server.
	NATSURL string `envconfig:"NATS_URL"`

	// MongoURI is the connection URI for the MongoDB server.
	MongoURI string `envconfig:"MONGO_URI"`

	// ServiceVersion is the version of the service, typically injected at build time.
	ServiceVersion string `envconfig:"SERVICE_VERSION" default:"dev"`
}

// New creates a new Config instance by loading values from environment variables.
// It returns an error if required environment variables are missing or invalid.
func New() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to process env vars: %w", err)
	}

	// Trim whitespace from all string fields.
	cfg.HTTPPort = strings.TrimSpace(cfg.HTTPPort)
	cfg.LogLevel = strings.TrimSpace(cfg.LogLevel)
	cfg.NATSURL = strings.TrimSpace(cfg.NATSURL)
	cfg.MongoURI = strings.TrimSpace(cfg.MongoURI)
	cfg.ServiceVersion = strings.TrimSpace(cfg.ServiceVersion)

	// NATS_URL is required; return a specific error if it's empty after trimming.
	if cfg.NATSURL == "" {
		return nil, ErrNATSURLRequired
	}

	return &cfg, nil
}
