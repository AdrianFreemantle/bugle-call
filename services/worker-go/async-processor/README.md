# Async Processor Microservice

A production-ready, modular Go microservice for asynchronous event processing in the Bugle Call platform. Exposes Prometheus metrics, health checks, and is designed for cloud-native, Kubernetes-based deployments.

## Features
- Modular, idiomatic Go codebase
- Environment-based configuration (see below)
- Structured JSON logging (slog)
- Prometheus `/metrics` endpoint
- `/healthz` readiness endpoint
- Graceful shutdown on SIGINT/SIGTERM
- Ready for NATS and MongoDB integration

## Configuration
Configuration is via environment variables (see also your Kubernetes manifests):

| Variable    | Required | Default | Description                 |
|-------------|----------|---------|-----------------------------|
| `NATS_URL`  | Yes      |         | NATS connection URL         |
| `HTTP_PORT` | No       | 8080    | HTTP server port            |
| `LOG_LEVEL` | No       | info    | Log level: debug/info/warn/error |

## Endpoints
- `GET /healthz` — Health check, returns 200 OK
- `GET /metrics` — Prometheus metrics

## Running Locally
```sh
export NATS_URL=nats://localhost:4222
export HTTP_PORT=8080
export LOG_LEVEL=debug

go run ./cmd/async-processor
```

## Testing
Run all tests:
```sh
go test ./...
```

## Kubernetes
See `k8s/async-processor-deployment.yaml` for deployment example.

---

© Bugle Call Platform
