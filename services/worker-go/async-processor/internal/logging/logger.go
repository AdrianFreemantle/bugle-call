// Logger setup for Async Processor
// Provides a structured slog-based logger.

package logging

import (
  "log/slog"
  "os"
)

// NewLogger returns a configured slog.Logger.
func NewLogger(level string) *slog.Logger {
  var lvl slog.Level
  switch level {
  case "DEBUG":
    lvl = slog.LevelDebug
  case "INFO":
    lvl = slog.LevelInfo
  case "WARN":
    lvl = slog.LevelWarn
  case "ERROR":
    lvl = slog.LevelError
  default:
    lvl = slog.LevelInfo
  }
  handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
  return slog.New(handler)
}
