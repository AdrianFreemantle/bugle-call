// Test for config loading
package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_Success(t *testing.T) {
	t.Setenv("NATS_URL", "nats://localhost:4222")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("HTTP_PORT", "9090")
	t.Setenv("MONGO_URI", "mongodb://localhost:27017")

	cfg, err := New()
	require.NoError(t, err)

	assert.Equal(t, "nats://localhost:4222", cfg.NATSURL)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, "9090", cfg.HTTPPort)
	assert.Equal(t, "mongodb://localhost:27017", cfg.MongoURI)
}

func TestNew_MissingRequired(t *testing.T) {
	t.Run("missing NATS_URL triggers error", func(t *testing.T) {
		t.Setenv("NATS_URL", "")
		_, err := New()
		if err == nil {
			t.Error("expected error when NATS_URL is missing")
		}
		if err != nil && !strings.Contains(err.Error(), "NATS_URL") {
			t.Errorf("expected error to mention NATS_URL, got: %v", err)
		}
	})

	t.Run("NATS_URL is single space triggers error", func(t *testing.T) {
		t.Setenv("NATS_URL", " ")
		_, err := New()
		if err == nil {
			t.Error("expected error when NATS_URL is a single space")
		}
		if err != nil && !strings.Contains(err.Error(), "NATS_URL") {
			t.Errorf("expected error to mention NATS_URL, got: %v", err)
		}
	})
}

func TestNew_InvalidNATSURL(t *testing.T) {
	t.Setenv("NATS_URL", "invalid-url")

	_, err := New()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid NATS_URL")
}

func TestNew_DefaultValues(t *testing.T) {
	// Save original environment
	originalEnv := os.Environ()
	
	// Clear all environment variables
	os.Clearenv()
	
	// Set only the required NATS_URL
	os.Setenv("NATS_URL", "nats://localhost:4222")
	
	// Create config with defaults
	cfg, err := New()
	require.NoError(t, err)
	
	// Restore original environment after test
	os.Clearenv()
	for _, envVar := range originalEnv {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 {
			os.Setenv(parts[0], parts[1])
		}
	}
	
	// Verify defaults
	assert.Equal(t, "nats://localhost:4222", cfg.NATSURL)
	assert.Equal(t, "info", cfg.LogLevel, "LogLevel should default to 'info'")
	assert.Equal(t, "8080", cfg.HTTPPort, "HTTPPort should default to '8080'")
}

// Mock for os.Exit to avoid terminating tests
func TestMustLoad_ExitsOnError(t *testing.T) {
	// Create a temporary test directory
	tempDir := t.TempDir()
	
	// Save original environment and functions
	originalEnv := os.Environ()
	originalOsExit := osExit
	originalStderr := os.Stderr
	
	// Create a temporary stderr file
	stderrFile := tempDir + "/stderr.txt"
	f, err := os.Create(stderrFile)
	require.NoError(t, err)
	
	// Set up test environment
	os.Clearenv()
	os.Setenv("NATS_URL", "") // Invalid empty URL
	os.Stderr = f
	
	// Mock os.Exit
	exitCalled := false
	exitCode := 0
	osExit = func(code int) {
		exitCalled = true
		exitCode = code
	}
	
	// Ensure environment is restored after test
	defer func() {
		// Restore original environment
		os.Clearenv()
		for _, envVar := range originalEnv {
			parts := strings.SplitN(envVar, "=", 2)
			if len(parts) == 2 {
				os.Setenv(parts[0], parts[1])
			}
		}
		
		// Restore original functions and stderr
		osExit = originalOsExit
		os.Stderr = originalStderr
		f.Close()
	}()
	
	// Call MustLoad (should call our mocked os.Exit)
	MustLoad()
	
	// Verify that os.Exit was called with code 1
	assert.True(t, exitCalled, "os.Exit should have been called")
	assert.Equal(t, 1, exitCode, "os.Exit should have been called with code 1")
	
	// Verify that error message was written to stderr
	f.Sync()
	f.Seek(0, 0)
	stderrBytes, err := os.ReadFile(stderrFile)
	require.NoError(t, err)
	assert.Contains(t, string(stderrBytes), "FATAL: invalid NATS_URL")
}
