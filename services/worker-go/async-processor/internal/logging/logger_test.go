package logging

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestNewLogger_LogLevelParsing(t *testing.T) {
	cases := []struct {
		level    string
		expected slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"error", slog.LevelError},
		{"invalid", slog.LevelInfo}, // fallback
		{"", slog.LevelInfo},        // fallback
	}
	for _, c := range cases {
		buf := &bytes.Buffer{}
		logger := NewLoggerWithWriter(c.level, buf)
		logger.Debug("debug msg")
		logger.Info("info msg")
		logger.Warn("warn msg")
		logger.Error("error msg")
		logOutput := buf.String()
		if c.expected == slog.LevelDebug && !strings.Contains(logOutput, "debug msg") {
			t.Errorf("expected debug msg for level %s", c.level)
		}
		if c.expected > slog.LevelDebug && strings.Contains(logOutput, "debug msg") {
			t.Errorf("did not expect debug msg for level %s", c.level)
		}
	}
}

func TestLogger_GlobalServiceField(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLoggerWithWriter("info", buf)
	logger.Info("test message")
	var m map[string]any
	err := json.Unmarshal(buf.Bytes(), &m)
	if err != nil {
		t.Fatalf("failed to unmarshal log: %v", err)
	}
	if m["service"] != "async-processor" {
		t.Errorf("expected service field to be async-processor, got %v", m["service"])
	}
}

func TestLogStartupBanner_OutputsExpectedFields(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLoggerWithWriter("info", buf)
	LogStartupBanner(logger, "async-processor", "v1.2.3", "nats://test", "8080")
	var m map[string]any
	err := json.Unmarshal(buf.Bytes(), &m)
	if err != nil {
		t.Fatalf("failed to unmarshal log: %v", err)
	}
	if m["msg"] != "starting service" {
		t.Errorf("expected msg field to be 'starting service', got %v", m["msg"])
	}
	if m["service"] != "async-processor" {
		t.Errorf("expected service field to be async-processor, got %v", m["service"])
	}
	if m["version"] != "v1.2.3" || m["nats_url"] != "nats://test" || m["http_port"] != "8080" {
		t.Errorf("unexpected banner fields: %+v", m)
	}
}

// NewLoggerWithWriter is a test helper to create a logger that writes to a buffer.
func NewLoggerWithWriter(level string, w *bytes.Buffer) *slog.Logger {
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
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: lvl,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				return slog.Attr{} // Remove source if present
			}
			return a
		},
	})
	return slog.New(handler).With("service", "async-processor")
}
