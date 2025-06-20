// HTTP server setup for Async Processor
// Exposes /metrics endpoint using Prometheus promhttp.Handler().

package http

import (
  "context"
  "fmt"
  "net/http"
  "os"
  "os/signal"
  "syscall"
  "time"

  "github.com/prometheus/client_golang/prometheus/promhttp"
  "log/slog"
)

// Server wraps the HTTP server logic.
type Server struct {
  srv    *http.Server
  logger *slog.Logger
}

// NewServer creates a new HTTP server with /metrics endpoint.
func NewServer(addr string, logger *slog.Logger) *Server {
  mux := http.NewServeMux()
  // /metrics endpoint for Prometheus
  mux.Handle("/metrics", promhttp.Handler())

  srv := &http.Server{
    Addr:    addr,
    Handler: mux,
  }
  return &Server{srv: srv, logger: logger}
}

// Start runs the HTTP server in a goroutine and handles graceful shutdown.
func (s *Server) Start(ctx context.Context) error {
  errCh := make(chan error, 1)
  go func() {
    s.logger.Info("starting HTTP server", "addr", s.srv.Addr)
    errCh <- s.srv.ListenAndServe()
  }()

  // Listen for shutdown signals
  stop := make(chan os.Signal, 1)
  signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

  select {
  case <-ctx.Done():
    s.logger.Info("context cancelled, shutting down HTTP server")
  case sig := <-stop:
    s.logger.Info("received signal, shutting down HTTP server", "signal", sig)
  case err := <-errCh:
    if err != nil && err != http.ErrServerClosed {
      s.logger.Error("HTTP server error", "err", err)
      return err
    }
    return nil
  }

  // Graceful shutdown
  ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  return s.srv.Shutdown(ctxTimeout)
}
