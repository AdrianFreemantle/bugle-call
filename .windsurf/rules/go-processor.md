---
trigger: glob
globs: src/processor/**/*.go
---

# Linting and Formatting
Use gofumpt for formatting (not just gofmt)
Enforce linting via golangci-lint with errcheck, govet, revive

# Package Layout
Use cmd/, internal/, and pkg/ folder convention
Do not mix app startup with processing logic
Keep NATS logic abstracted under internal/messaging/

# Message Handling
Use context.Context on all message handlers
Acknowledge NATS messages manually after successful handling
Implement retry logic with exponential backoff
Log every failed message with trace context

# Metrics
Expose /metrics using Prometheus Go client
Avoid global metric definitions
Label metrics with queue name and message type