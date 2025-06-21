// Entry point for Async Processor microservice
// Loads config, initializes logger, starts HTTP server, and fails fast on errors.

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/example/async-processor/internal/config"
	apphttp "github.com/example/async-processor/internal/http"
	"github.com/example/async-processor/internal/logging"
	"github.com/example/async-processor/internal/subscriber"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "async-processor error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Create a context that is canceled on SIGINT or SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Load configuration
	cfg, err := config.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger := logging.NewLogger(cfg.LogLevel, "async-processor", os.Stdout)
	version := os.Getenv("SERVICE_VERSION")
	if version == "" {
		version = "dev"
	}
	logging.LogStartupBanner(logger, version, cfg.NATSURL, cfg.HTTPPort)

	// Start components
	httpSrv := apphttp.NewServer(cfg.HTTPPort)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpSrv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server error", "err", err)
			stop() // trigger shutdown
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := subscriber.Start(ctx, cfg, logger); err != nil && !errors.Is(err, context.Canceled) {
			logger.Error("subscriber error", "err", err)
			stop() // trigger shutdown
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()

	// Initiate graceful shutdown
	logger.Info("shutdown signal received, starting graceful shutdown")

	// Create a context with a timeout for shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		logger.Error("http server shutdown error", "err", err)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	logger.Info("shutdown complete")

	return nil
}
