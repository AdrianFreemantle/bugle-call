package http

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func startTestServer(t *testing.T) (*Server, string, func()) {
	// Listen on random port
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	addr := l.Addr().String()
	srv := NewServer()
	srv.httpServer.Addr = addr

	_, cancel := context.WithCancel(context.Background())
	go func() {
		_ = srv.httpServer.Serve(l)
	}()

	// Wait for server to be ready
	time.Sleep(100 * time.Millisecond)
	return srv, addr, func() {
		cancel()
		ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancelTimeout()
		_ = srv.httpServer.Shutdown(ctxTimeout)
	}
}

func TestHealthzEndpoint(t *testing.T) {
	_, addr, shutdown := startTestServer(t)
	defer shutdown()

	resp, err := http.Get("http://" + addr + "/healthz")
	require.NoError(t, err, "failed to GET /healthz")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "expected 200 OK")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read response body")
	assert.Equal(t, "ok", string(body), "expected body to be 'ok'")
}

func TestMetricsEndpoint(t *testing.T) {
	_, addr, shutdown := startTestServer(t)
	defer shutdown()

	resp, err := http.Get("http://" + addr + "/metrics")
	if err != nil {
		t.Fatalf("failed to GET /metrics: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "go_gc_duration_seconds") {
		t.Errorf("expected Prometheus metrics in response, got %q", body[:100])
	}
}

func TestGracefulShutdown(t *testing.T) {
	_, addr, shutdown := startTestServer(t)

	// Make a request before shutdown
	resp, err := http.Get("http://" + addr + "/healthz")
	if err != nil {
		t.Fatalf("failed to GET /healthz: %v", err)
	}
	resp.Body.Close()

	// Shutdown server
	shutdown()

	// Wait for shutdown to complete
	time.Sleep(100 * time.Millisecond)

	// Next request should fail
	_, err = http.Get("http://" + addr + "/healthz")
	if err == nil {
		t.Error("expected error after shutdown, got nil")
	}
}
