// HTTP server setup for Async Processor
// Exposes /metrics endpoint using Prometheus promhttp.Handler().

// Package httpserver provides HTTP server functionality for the async processor.
// It exposes /metrics and /healthz endpoints, supports Prometheus, and allows graceful shutdown.
package http

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server is a production-grade HTTP server for the async processor.
type Server struct {
	httpServer *http.Server
}

// NewServer creates a new HTTP server listening on :8080 with /metrics and /healthz endpoints.
func NewServer() *Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	httpSrv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	return &Server{httpServer: httpSrv}
}

// Start runs the HTTP server and blocks until the context is cancelled or server is shut down.
func (s *Server) Start(ctx context.Context) error {
	done := make(chan error, 1)
	go func() {
		done <- s.httpServer.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-done:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	}
}

// Stop gracefully shuts down the HTTP server with the given context timeout
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
