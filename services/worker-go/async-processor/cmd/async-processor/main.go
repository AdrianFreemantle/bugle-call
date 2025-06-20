// Entry point for Async Processor microservice
// Loads config, initializes logger, starts HTTP server, and fails fast on errors.

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Create error channels for components
	httpErrCh := make(chan error, 1)
	subErrCh := make(chan error, 1)

	// Start HTTP server
	httpSrv := http.NewServer()
	go func() {
		httpErrCh <- httpSrv.Start(ctx)
	}()

	// Start NATS subscriber
	go func() {
		if err := subscriber.Start(ctx, cfg, logger); err != nil {
			subErrCh <- err
		}
	}()

	// Wait for termination or error
	var exitCode int
	select {
	case err := <-httpErrCh:
		if err != nil {
			logger.Error("HTTP server exited with error", "err", err)
			exitCode = 1
		}
	case err := <-subErrCh:
		logger.Error("Subscriber exited with error", "err", err)
		exitCode = 1
	case <-ctx.Done():
		logger.Info("received shutdown signal")
	}

	// Initiate graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown HTTP server
	if err := httpSrv.Stop(shutdownCtx); err != nil {
		logger.Error("error shutting down HTTP server", "err", err)
		exitCode = 1
	}

	logger.Info("async-processor shutdown complete")
	if exitCode != 0 {
		os.Exit(exitCode)
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
