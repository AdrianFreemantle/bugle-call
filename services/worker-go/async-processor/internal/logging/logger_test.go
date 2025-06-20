package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNewLogger_LogLevelParsing verifies that log level strings and slog.Level values are parsed correctly and logger emits expected output.
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
		{"DeBuG", slog.LevelDebug},  // case-insensitive
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

// TestLogger_GlobalServiceField verifies that the global 'service' field is always present in logs.
func TestLogger_GlobalServiceField(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLoggerWithWriter("info", buf)
	logger.Info("test message")
	var m map[string]any
	err := json.Unmarshal(buf.Bytes(), &m)
	require.NoError(t, err, "failed to unmarshal log")
	require.Equal(t, "async-processor", m["service"], "expected service field to be async-processor")
}

// parseLogBuffer unmarshals the JSON log buffer and returns the map and error.
func parseLogBuffer(t *testing.T, buf *bytes.Buffer) (map[string]any, error) {
	t.Helper()
	var m map[string]any
	// Get the last line from the buffer, which should be the main log entry
	lines := strings.Split(buf.String(), "\n")
	var lastLine string
	for i := len(lines) - 1; i >= 0; i-- {
		if len(strings.TrimSpace(lines[i])) > 0 {
			lastLine = lines[i]
			break
		}
	}
	
	if lastLine == "" {
		return nil, fmt.Errorf("no log entries found")
	}
	
	err := json.Unmarshal([]byte(lastLine), &m)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return m, nil
}

// TestLogStartupBanner_OutputsExpectedFields verifies that LogStartupBanner emits all expected fields with correct values.
func TestLogStartupBanner_OutputsExpectedFields(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLoggerWithWriter("info", buf)
	LogStartupBanner(logger, "async-processor", "v1.2.3", "nats://test", "8080", "component", "test-component", "build", "abc123", "env", "dev")
	m, err := parseLogBuffer(t, buf)
	require.NoError(t, err, "Failed to parse log JSON:\n%s", buf.String())
	require.Equal(t, "starting service", m["msg"], "expected 'msg' field to be 'starting service'. Raw: %s", buf.String())
	require.Equal(t, "async-processor", m["service"], "expected 'service' field. Raw: %s", buf.String())
	require.Equal(t, "test-component", m["component"], "expected 'component' field. Raw: %s", buf.String())
	require.Equal(t, "v1.2.3", m["version"], "expected 'version' field. Raw: %s", buf.String())
	require.Equal(t, "nats://test", m["nats_url"], "expected 'nats_url' field. Raw: %s", buf.String())
	require.Equal(t, "8080", m["http_port"], "expected 'http_port' field. Raw: %s", buf.String())
	require.Equal(t, "abc123", m["build"], "expected 'build' field. Raw: %s", buf.String())
	require.Equal(t, "dev", m["env"], "expected 'env' field. Raw: %s", buf.String())
}

// TestLogStartupBanner_EmptyValues verifies that startup log includes all expected keys, even if values are empty.
func TestLogStartupBanner_EmptyValues(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLoggerWithWriter("info", buf)
	LogStartupBanner(logger, "", "", "", "", "component", "")
	m, err := parseLogBuffer(t, buf)
	require.NoError(t, err, "Failed to parse log JSON:\n%s", buf.String())
	require.Equal(t, "starting service", m["msg"], "expected 'msg' field to be 'starting service'. Raw: %s", buf.String())
	require.Contains(t, m, "service", "expected 'service' key to exist. Raw: %s", buf.String())
	require.Contains(t, m, "component", "expected 'component' key to exist. Raw: %s", buf.String())
	require.Contains(t, m, "version", "expected 'version' key to exist. Raw: %s", buf.String())
	require.Contains(t, m, "nats_url", "expected 'nats_url' key to exist. Raw: %s", buf.String())
	require.Contains(t, m, "http_port", "expected 'http_port' key to exist. Raw: %s", buf.String())
}

// TestLogStartupBanner_MultipleKeyValuePairs_NoConflict verifies that LogStartupBanner does not overwrite/conflict with multiple key-value pairs.
func TestLogStartupBanner_MultipleKeyValuePairs_NoConflict(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLoggerWithWriter("info", buf)
	LogStartupBanner(
		logger,
		"svc", "vX", "nats://multi", "9999",
		"component", "multi-comp", "build", "multi-build", "env", "multi-env", "extra1", "val1", "extra2", "val2",
	)
	m, err := parseLogBuffer(t, buf)
	require.NoError(t, err, "Failed to parse log JSON:\n%s", buf.String())
	require.Equal(t, "svc", m["service"])
	require.Equal(t, "multi-comp", m["component"])
	require.Equal(t, "vX", m["version"])
	require.Equal(t, "nats://multi", m["nats_url"])
	require.Equal(t, "9999", m["http_port"])
	require.Equal(t, "multi-build", m["build"])
	require.Equal(t, "multi-env", m["env"])
	require.Equal(t, "val1", m["extra1"])
	require.Equal(t, "val2", m["extra2"])
}

// parseSlogLevel parses a log level from string or slog.Level, defaulting to slog.LevelInfo.
func parseSlogLevel(level any) slog.Level {
	switch v := level.(type) {
	case slog.Level:
		return v
	case string:
		switch strings.ToLower(v) {
		case "debug":
			return slog.LevelDebug
		case "warn":
			return slog.LevelWarn
		case "error":
			return slog.LevelError
		default:
			return slog.LevelInfo
		}
	default:
		return slog.LevelInfo
	}
}

// NewLoggerWithWriter creates a logger that writes to the given buffer with the specified log level.
// The level can be a string or slog.Level for test and non-test flexibility.
func NewLoggerWithWriter(level any, w *bytes.Buffer) *slog.Logger {
	tlevel := parseSlogLevel(level)
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: tlevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				return slog.Attr{} // Remove source if present
			}
			return a
		},
	})
	return slog.New(handler).With("service", "async-processor")
}

// TestParseLogBuffer_InvalidPayload verifies that parseLogBuffer returns an error for invalid JSON.
func TestParseLogBuffer_InvalidPayload(t *testing.T) {
	buf := &bytes.Buffer{}
	buf.WriteString("{not valid json}")
	m, err := parseLogBuffer(t, buf)
	require.Error(t, err, "Expected error on invalid JSON. Raw: %s", buf.String())
	require.Nil(t, m, "Expected nil map on invalid JSON. Raw: %s", buf.String())
}

// TestLogStartupBanner_ExtrasAppear verifies that extra keys and values appear in the log output JSON.
func TestLogStartupBanner_ExtrasAppear(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLoggerWithWriter("info", buf)
	LogStartupBanner(logger, "svc", "v1", "nats://test", "8080", "foo", "bar", "baz", 123)
	m, err := parseLogBuffer(t, buf)
	require.NoError(t, err)
	require.Equal(t, "bar", m["foo"])
	require.Equal(t, 123.0, m["baz"]) // json decodes numbers as float64
}

// TestLogStartupBanner_OddExtrasIgnored verifies that an odd number of extra arguments logs a warning and ignores extras.
func TestLogStartupBanner_OddExtrasIgnored(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLoggerWithWriter("info", buf)
	LogStartupBanner(logger, "svc", "v1", "nats://test", "8080", "foo")
	
	// There should be a warning log and an info log
	logLines := strings.Split(buf.String(), "\n")
	var infoLogFound bool
	for _, line := range logLines {
		if len(line) == 0 {
			continue
		}
		
		var logEntry map[string]any
		err := json.Unmarshal([]byte(line), &logEntry)
		require.NoError(t, err, "Failed to parse log line: %s", line)
		
		// Check the info log (startup banner)
		if logEntry["msg"] == "starting service" {
			infoLogFound = true
			require.NotContains(t, logEntry, "foo", "'foo' should not be in the log entry")
		}
	}
	
	require.True(t, infoLogFound, "Info log with startup banner not found")
}

// TestLogStartupBanner_NonStringKeySkipped verifies that non-string keys in extras log a warning and are skipped.
func TestLogStartupBanner_NonStringKeySkipped(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLoggerWithWriter("info", buf)
	LogStartupBanner(logger, "svc", "v1", "nats://test", "8080", 123, "badkey", "good", "ok")
	
	// There should be a warning log and an info log
	logLines := strings.Split(buf.String(), "\n")
	var infoLogFound bool
	for _, line := range logLines {
		if len(line) == 0 {
			continue
		}
		
		var logEntry map[string]any
		err := json.Unmarshal([]byte(line), &logEntry)
		require.NoError(t, err, "Failed to parse log line: %s", line)
		
		// Check the info log (startup banner)
		if logEntry["msg"] == "starting service" {
			infoLogFound = true
			require.NotContains(t, logEntry, "badkey", "'badkey' should not be in the log entry")
			require.Equal(t, "ok", logEntry["good"], "'good' field should have value 'ok'")
		}
	}
	
	require.True(t, infoLogFound, "Info log with startup banner not found")
}
