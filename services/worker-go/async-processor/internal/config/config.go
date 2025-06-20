// Config loading logic for Async Processor
// Loads configuration from environment variables using envconfig.
// Fails fast if any required config is missing or invalid.

package config

import (
  "fmt"
  "os"
  "github.com/kelseyhightower/envconfig"
)

// Config holds all configuration for the async processor.
type Config struct {
  HTTPPort string `envconfig:"HTTP_PORT" required:"true"`
  LogLevel string `envconfig:"LOG_LEVEL" default:"INFO"`
}

// Load loads config from env vars, fails fast on error.
func Load() (*Config, error) {
  var cfg Config
  if err := envconfig.Process("", &cfg); err != nil {
    return nil, fmt.Errorf("failed to load config: %w", err)
  }
  return &cfg, nil
}

// MustLoad loads config and exits on error (fail-fast).
func MustLoad() *Config {
  cfg, err := Load()
  if err != nil {
    fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
    os.Exit(1)
  }
  return cfg
}
