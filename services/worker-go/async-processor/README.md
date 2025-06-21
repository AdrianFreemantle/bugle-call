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

Configuration is managed via environment variables.

### Required

| Variable   | Description               |
|------------|---------------------------|
| `NATS_URL` | NATS server connection URL. |

### Optional

| Variable        | Default | Description                                  |
|-----------------|---------|----------------------------------------------|
| `HTTP_PORT`     | `8080`  | Port for the HTTP server.                    |
| `LOG_LEVEL`     | `info`  | Log level (`debug`, `info`, `warn`, `error`).  |
| `MONGO_URI`     | `""`    | MongoDB connection URI (for future use).     |
| `SERVICE_VERSION` | `dev`   | Service version, included in startup logs.   |

## API Endpoints

- `GET /healthz`: Health check endpoint. Returns `200 OK` with body `ok`.
- `GET /metrics`: Exposes Prometheus metrics.

## Running Locally

1.  **Set Environment Variables**:

    ```sh
    export NATS_URL="nats://localhost:4222"
    export HTTP_PORT="8080"
    export LOG_LEVEL="debug"
    ```

2.  **Run the Service**:

    ```sh
    go run ./cmd/async-processor
    ```

3.  **Test Endpoints**:

    ```sh
    # Check health
    curl http://localhost:8080/healthz

    # View metrics
    curl http://localhost:8080/metrics
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
