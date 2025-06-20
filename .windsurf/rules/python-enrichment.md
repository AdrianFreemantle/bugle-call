---
trigger: glob
globs: src/enrichment/**/*.py
---

# Linting and Formatting
Use black for formatting
Use ruff for linting
All function signatures must include type hints
Strict mode required (pyproject.toml: strict = true)

# Structure and Model Hygiene
Separate rule-based logic from ML model code
Keep all enrichment logic inside enrichment/ or pipeline/
Include a version string in the enrichment response
Models must run locally without calling external APIs

# API Design
Use FastAPI for consistent structure
Expose /metrics (Prometheus format) and /healthz endpoints
Use structured logging with trace ID per request
Do not log full payloads in production logs

# Confidence and Traceability
Return classifier confidence in every response
Log enrichment duration and model version
Support trace IDs passed in headers (e.g., X-Request-ID)
