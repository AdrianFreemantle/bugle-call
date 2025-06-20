// Entry point for Async Processor microservice
// Loads config, initializes logger, starts HTTP server, and fails fast on errors.

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/example/async-processor/internal/config"
	"github.com/example/async-processor/internal/http"
	"github.com/example/async-processor/internal/logging"
	"github.com/example/async-processor/internal/subscriber"
)

func main() {
	// Load configuration
	cfg := config.MustLoad()

	// Initialize logger
	logger := logging.NewLogger(cfg.LogLevel)
	version := os.Getenv("SERVICE_VERSION")
	if version == "" {
		version = "dev"
	}
	logging.LogStartupBanner(logger, "async-processor", version, cfg.NATSURL, cfg.HTTPPort)

	// Set up context with signal handling for graceful shutdown
	ctx, stop := signalContext()
	defer stop()

	// Start HTTP server
	httpSrv := http.NewServer()
	httpErrCh := make(chan error, 1)
	go func() {
		httpErrCh <- httpSrv.Start(ctx)
	}()

	// Start NATS subscriber stub
	go func() {
		_ = subscriber.Start(ctx, cfg, logger)
	}()

	// Wait for termination or error
	select {
	case err := <-httpErrCh:
		if err != nil {
			logger.Error("HTTP server exited with error", "err", err)
			os.Exit(1)
		}
	case <-ctx.Done():
		logger.Info("shutting down async-processor")
	}
}

// signalContext returns a context that is cancelled on SIGINT or SIGTERM
func signalContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		cancel()
	}()
	return ctx, cancel
}
