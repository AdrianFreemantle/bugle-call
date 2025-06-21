package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_Endpoints(t *testing.T) {
	// Find a free port and create a server instance
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	addr := l.Addr().String()
	port := fmt.Sprintf("%d", l.Addr().(*net.TCPAddr).Port)
	l.Close() // Close the listener, as Start() will create its own.

	server := NewServer(port)

	// Start the server in a goroutine so it doesn't block
	go func() {
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("HTTP server ListenAndServe: %v", err)
		}
	}()

	// Defer the shutdown to ensure the server is cleaned up
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		require.NoError(t, server.Shutdown(shutdownCtx))
	}()

	// Wait until the server is ready to accept connections
	require.Eventually(t, func() bool {
		conn, err := net.DialTimeout("tcp", addr, 50*time.Millisecond)
		if err != nil {
			return false
		}
		conn.Close()
		return true
	}, 2*time.Second, 100*time.Millisecond, "server was not ready")

	t.Run("/healthz", func(t *testing.T) {
		resp, err := http.Get("http://" + addr + "/healthz")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, "ok\n", string(body))
	})

	t.Run("/metrics", func(t *testing.T) {
		resp, err := http.Get("http://" + addr + "/metrics")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Contains(t, string(body), "go_gc_duration_seconds")
	})
}

func TestServer_Shutdown(t *testing.T) {
	// Find a free port and create a server instance
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	addr := l.Addr().String()
	port := fmt.Sprintf("%d", l.Addr().(*net.TCPAddr).Port)
	l.Close() // Close the listener, as Start() will create its own.

	srv := NewServer(port)

	// Run the server in a goroutine so it doesn't block
	go func() {
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("HTTP server ListenAndServe: %v", err)
		}
	}()

	// Wait until the server is ready to accept connections
	require.Eventually(t, func() bool {
		conn, err := net.DialTimeout("tcp", addr, 50*time.Millisecond)
		if err != nil {
			return false
		}
		conn.Close()
		return true
	}, 2*time.Second, 100*time.Millisecond, "server was not ready")

	// Now, initiate a graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = srv.Shutdown(shutdownCtx)
	require.NoError(t, err, "expected clean shutdown")

	// Verify the server is no longer listening
	conn, err := net.DialTimeout("tcp", addr, 200*time.Millisecond)
	assert.Error(t, err, "server should not be listening after shutdown")
	if conn != nil {
		conn.Close()
	}
}
