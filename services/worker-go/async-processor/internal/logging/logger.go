// Logger setup for Async Processor
// Provides a structured slog-based logger.

// Package logging provides structured, production-grade logging via slog.
// It supports JSON output and environment-configurable log levels.
package logging

import (
	"io"
	"log/slog"
	"strings"
)

// NewLogger returns a slog.Logger with JSON output and the specified log level.
// It parses the level case-insensitively and defaults to "info".
func NewLogger(level string, serviceName string, w io.Writer) *slog.Logger {
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo // Default to info for invalid or empty levels
	}

	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{Level: lvl})
	return slog.New(handler).With("service", serviceName)
}

// LogStartupBanner logs a service startup banner with key config fields and any optional extra fields.
func LogStartupBanner(logger *slog.Logger, version, natsURL, httpPort string, extra ...any) {
	args := []any{
		slog.String("version", version),
		slog.String("nats_url", natsURL),
		slog.String("http_port", httpPort),
	}
	args = append(args, extra...)

	logger.Info("starting service", args...)
}
