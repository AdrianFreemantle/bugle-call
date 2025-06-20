// Test for HTTP server
package http

import (
  "net/http"
  "net/http/httptest"
  "testing"
  "github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestMetricsHandler(t *testing.T) {
  mux := http.NewServeMux()
  mux.Handle("/metrics", promhttp.Handler())
  srv := httptest.NewServer(mux)
  defer srv.Close()

  resp, err := http.Get(srv.URL + "/metrics")
  if err != nil {
    t.Fatalf("failed to GET /metrics: %v", err)
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    t.Errorf("expected status 200 OK, got %d", resp.StatusCode)
  }
}
