// Logger setup for Async Processor
// Provides a structured slog-based logger.

// Package logging provides structured, production-grade logging via slog.
// Supports JSON output and environment-configurable log levels.
package logging

import (
	"log/slog"
	"os"
	"strings"
)

// NewLogger returns a slog.Logger with JSON output and the specified log level.
// Level should be one of: "debug", "info", "warn", "error" (case-insensitive).
func NewLogger(level string) *slog.Logger {
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	return slog.New(handler)
}

// LogStartupBanner logs a service startup banner with key config fields.
func LogStartupBanner(logger *slog.Logger, service, version, natsURL, httpPort string) {
	logger.Info("starting service",
		"service", service,
		"version", version,
		"nats_url", natsURL,
		"http_port", httpPort,
	)
}
