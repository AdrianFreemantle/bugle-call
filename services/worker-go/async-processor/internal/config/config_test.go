// Test for config loading
package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_Success(t *testing.T) {
	t.Setenv("NATS_URL", "nats://localhost:4222")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("HTTP_PORT", "9090")

	cfg, err := New()
	require.NoError(t, err)

	assert.Equal(t, "nats://localhost:4222", cfg.NATSURL)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, "9090", cfg.HTTPPort)
}

func TestNew_MissingRequired(t *testing.T) {
	// Clear NATS_URL which is required
	t.Setenv("NATS_URL", "")

	_, err := New()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "NATS_URL")
}

func TestNew_InvalidNATSURL(t *testing.T) {
	t.Setenv("NATS_URL", "invalid-url")

	_, err := New()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid NATS_URL")
}

func TestNew_DefaultValues(t *testing.T) {
	t.Setenv("NATS_URL", "nats://localhost:4222")
	// Clear optional fields to test defaults
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("HTTP_PORT", "")

	cfg, err := New()
	require.NoError(t, err)

	assert.Equal(t, "nats://localhost:4222", cfg.NATSURL)
	assert.Equal(t, "info", cfg.LogLevel)    // default
	assert.Equal(t, "8080", cfg.HTTPPort)    // default
}

func TestMustLoad_PanicsOnError(t *testing.T) {
	// Clear required NATS_URL to cause an error
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"test"}
	t.Setenv("NATS_URL", "") // Invalid empty URL

	// Redirect stderr to avoid test output pollution
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	done := make(chan bool)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- true
			}
		}()
		MustLoad()
		t.Error("MustLoad should have panicked")
	}()

	w.Close()
	os.Stderr = oldStderr

	// Wait for the panic or timeout
	select {
	case <-done:
		// Test passed - MustLoad panicked as expected
	case <-time.After(1 * time.Second):
		t.Error("MustLoad should have panicked")
	}
}
