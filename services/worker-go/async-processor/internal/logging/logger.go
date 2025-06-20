// Logger setup for Async Processor
// Provides a structured slog-based logger.

// Package logging provides structured, production-grade logging via slog.
// Supports JSON output and environment-configurable log levels.
package logging

import (
	"fmt"
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
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				return slog.Attr{} // Remove source if present
			}
			return a
		},
	})
	// Add global service field to all logs
	logger := slog.New(handler).With("service", "async-processor")
	return logger
}

// LogStartupBanner logs a service startup banner with key config fields and any optional extra fields.
// Example: LogStartupBanner(logger, "async-processor", "v1.2.3", "nats://test", "8080", "build", "abc123", "env", "dev")
func LogStartupBanner(logger *slog.Logger, service, version, natsURL, httpPort string, extra ...any) {
	attrs := []any{
		"service", service,
		"version", version,
		"nats_url", natsURL,
		"http_port", httpPort,
	}

	if len(extra)%2 != 0 {
		logger.Warn("LogStartupBanner: odd number of extra arguments; extras ignored", "extra_len", len(extra))
		logger.Info("starting service", attrs...)
		return
	}

	for i := 0; i < len(extra); i += 2 {
		key, ok := extra[i].(string)
		if !ok {
			logger.Warn("LogStartupBanner: non-string key in extra fields; pair skipped", "key_type", fmt.Sprintf("%T", extra[i]))
			continue
		}
		attrs = append(attrs, key, extra[i+1])
	}
	logger.Info("starting service", attrs...)
}
