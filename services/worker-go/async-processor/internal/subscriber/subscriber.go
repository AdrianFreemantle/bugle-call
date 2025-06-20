// Package subscriber provides the NATS event subscription logic for the async processor.
// This is a stub to prepare for NATS integration.
package subscriber

import (
	"context"
	"log/slog"

	"github.com/example/async-processor/internal/config"
)

// Start runs the NATS subscriber. This is a stub for future implementation.
func Start(ctx context.Context, cfg *config.Config, logger *slog.Logger) error {
	logger.Info("subscriber stub: NATS integration not yet implemented")
	<-ctx.Done()
	return ctx.Err()
}
