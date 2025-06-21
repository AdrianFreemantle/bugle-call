package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_Success(t *testing.T) {
	t.Setenv("NATS_URL", "nats://localhost:4222")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("HTTP_PORT", "9090")
	t.Setenv("MONGO_URI", "mongodb://user:pass@host:27017/db")
	t.Setenv("SERVICE_VERSION", "v1.2.3")

	cfg, err := New()
	require.NoError(t, err)

	assert.Equal(t, "nats://localhost:4222", cfg.NATSURL)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, "9090", cfg.HTTPPort)
	assert.Equal(t, "mongodb://user:pass@host:27017/db", cfg.MongoURI)
	assert.Equal(t, "v1.2.3", cfg.ServiceVersion)
}

func TestNew_RequiredFields(t *testing.T) {
	t.Run("NATS_URL is missing", func(t *testing.T) {
		t.Setenv("NATS_URL", "")
		_, err := New()
		require.Error(t, err)
		assert.EqualError(t, err, "NATS_URL is required and cannot be empty")
	})

	t.Run("NATS_URL is only whitespace", func(t *testing.T) {
		t.Setenv("NATS_URL", "   ")
		_, err := New()
		require.Error(t, err)
		assert.EqualError(t, err, "NATS_URL is required and cannot be empty")
	})
}

func TestNew_TrimsWhitespace(t *testing.T) {
	t.Setenv("NATS_URL", "  nats://localhost:4222  ")
	t.Setenv("LOG_LEVEL", "  debug  ")
	t.Setenv("HTTP_PORT", "  9090  ")
	t.Setenv("MONGO_URI", "  mongodb://localhost:27017  ")
	t.Setenv("SERVICE_VERSION", "  v1.0.0  ")

	cfg, err := New()
	require.NoError(t, err)

	assert.Equal(t, "nats://localhost:4222", cfg.NATSURL)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, "9090", cfg.HTTPPort)
	assert.Equal(t, "mongodb://localhost:27017", cfg.MongoURI)
	assert.Equal(t, "v1.0.0", cfg.ServiceVersion)
}

func TestNew_DefaultValues(t *testing.T) {
	t.Setenv("NATS_URL", "nats://default:4222")
	t.Setenv("MONGO_URI", "") // Explicitly unset

	cfg, err := New()
	require.NoError(t, err)

	assert.Equal(t, "8080", cfg.HTTPPort)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, "", cfg.MongoURI)
	assert.Equal(t, "nats://default:4222", cfg.NATSURL)
}
