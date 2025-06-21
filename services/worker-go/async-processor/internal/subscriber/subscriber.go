// Package subscriber provides the NATS event subscription logic for the async processor.
// This is a stub to prepare for NATS integration.
package subscriber

import (
	"context"
	"log/slog"

	"github.com/example/async-processor/internal/config"
)

// Start runs the NATS subscriber stub.
func Start(ctx context.Context, cfg *config.Config, logger *slog.Logger) error {
	logger.Info("NATS subscriber started", slog.String("nats_url", cfg.NATSURL))

	// Block until context is cancelled
	<-ctx.Done()

	logger.Info("NATS subscriber shutting down")
	return ctx.Err()
}
