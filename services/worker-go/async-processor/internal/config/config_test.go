// Test for config loading
package config

import (
  "os"
  "testing"
)

func TestLoad_Success(t *testing.T) {
  os.Setenv("HTTP_PORT", "8080")
  os.Setenv("LOG_LEVEL", "DEBUG")
  defer os.Unsetenv("HTTP_PORT")
  defer os.Unsetenv("LOG_LEVEL")

  cfg, err := Load()
  if err != nil {
    t.Fatalf("expected no error, got %v", err)
  }
  if cfg.HTTPPort != "8080" {
    t.Errorf("expected HTTPPort=8080, got %s", cfg.HTTPPort)
  }
  if cfg.LogLevel != "DEBUG" {
    t.Errorf("expected LogLevel=DEBUG, got %s", cfg.LogLevel)
  }
}

func TestLoad_MissingRequired(t *testing.T) {
  os.Unsetenv("HTTP_PORT")
  os.Unsetenv("LOG_LEVEL")
  _, err := Load()
  if err == nil {
    t.Fatal("expected error for missing HTTP_PORT, got nil")
  }
}
