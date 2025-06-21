package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestLogger creates a logger for testing purposes. It writes to a buffer
// and removes the timestamp to make log output deterministic for assertions.
func newTestLogger(level string, service string, w *bytes.Buffer) *slog.Logger {
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
			// Remove time for deterministic test output
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})
	return slog.New(handler).With("service", service)
}

// TestNewLogger tests that the NewLogger function correctly sets the log level.
func TestNewLogger(t *testing.T) {
	testCases := []struct {
		name          string
		level         string
		expectedLevel slog.Level
	}{
		{"LowercaseDebug", "debug", slog.LevelDebug},
		{"MixedCaseInfo", "InFo", slog.LevelInfo},
		{"UppercaseWarn", "WARN", slog.LevelWarn},
		{"InvalidLevel", "verbose", slog.LevelInfo},
		{"EmptyLevel", "", slog.LevelInfo},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Use the real NewLogger from the package.
			// We don't need to check output, just the handler's enabled state.
			logger := NewLogger(tc.level, "test-service", &bytes.Buffer{})

			assert.False(t, logger.Handler().Enabled(context.Background(), tc.expectedLevel-1))
			assert.True(t, logger.Handler().Enabled(context.Background(), tc.expectedLevel))
			assert.True(t, logger.Handler().Enabled(context.Background(), tc.expectedLevel+1))
		})
	}
}

// TestLogStartupBanner tests the output of the LogStartupBanner function.
func TestLogStartupBanner(t *testing.T) {
	var buf bytes.Buffer
	// Use the test-specific logger for deterministic output.
	logger := newTestLogger("info", "test-service", &buf)

	// Call the real LogStartupBanner from the package.
	LogStartupBanner(logger, "v1.2.3", "nats://test", "8080", "component", "tester")

	var output map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &output)
	require.NoError(t, err)

	assert.Equal(t, "starting service", output["msg"])
	assert.Equal(t, "test-service", output["service"])
	assert.Equal(t, "v1.2.3", output["version"])
	assert.Equal(t, "nats://test", output["nats_url"])
	assert.Equal(t, "8080", output["http_port"])
	assert.Equal(t, "tester", output["component"])
}

func TestLogStartupBanner_OddExtrasIgnored(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger("info", "test-service", &buf)

	LogStartupBanner(logger, "v1", "nats://test", "8080", "component", "tester", "foo")

	var output map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &output)
	require.NoError(t, err)

	assert.Equal(t, "starting service", output["msg"])
	assert.Equal(t, "tester", output["component"])
	_, exists := output["foo"]
	assert.False(t, exists, "unpaired attribute 'foo' should not be in log output")
}

func TestLogStartupBanner_NonStringKeySkipped(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger("info", "test-service", &buf)

	// We pass a valid pair before an invalid one to ensure slog processes
	// the valid attributes even when it encounters malformed ones.
	LogStartupBanner(logger, "v1", "nats://test", "8080", "good", "ok", 123, "badkey")

	var output map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &output)
	require.NoError(t, err)

	// The valid field should be present.
	assert.Equal(t, "ok", output["good"])

	// Slog should also create a '!BADKEY' field for the non-string key.
	_, badKeyExists := output["!BADKEY"]
	assert.True(t, badKeyExists, "slog should add a '!BADKEY' for non-string keys")
}
