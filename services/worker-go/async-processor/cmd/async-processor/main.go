// Entry point for Async Processor microservice
// Loads config, initializes logger, starts HTTP server, and fails fast on errors.

package main

import (
  "context"
  "fmt"
  "os"
  "time"

  "github.com/example/async-processor/internal/config"
  "github.com/example/async-processor/internal/http"
  "github.com/example/async-processor/internal/logging"
)

func main() {
  // Load configuration from environment variables
  cfg := config.MustLoad()

  // Initialize structured logger
  logger := logging.NewLogger(cfg.LogLevel)
  logger.Info("config loaded", "http_port", cfg.HTTPPort, "log_level", cfg.LogLevel)

  // Start HTTP server with /metrics endpoint
  srv := http.NewServer(":"+cfg.HTTPPort, logger)
  ctx := context.Background()
  if err := srv.Start(ctx); err != nil {
    logger.Error("failed to start HTTP server", "err", err)
    fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
    os.Exit(1)
  }

  // Block forever (or until shutdown)
  select {}
}
