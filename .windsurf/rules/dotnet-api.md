---
trigger: glob
globs: src/api/**/*.cs
---

# Language Features
Enable nullable reference types (<Nullable>enable</Nullable>)
Use record types for immutable DTOs
Avoid dynamic or object-based models in favor of strong typing

# Project Structure
Use PascalCase for class and file names
Keep domain models in Domain/ and transport models in DTOs/
Expose endpoints via Minimal API or clearly separated Controllers

# Security and Auth
Never hardcode secrets or signing keys
All JWT settings must come from environment or injected config
JWTs must be validated with audience, issuer, and expiry
Include signing algorithm (HS256 or RS256)

# Observability
Use prometheus-net.AspNetCore to expose /metrics endpoint
Tag metrics with request method, status code, and duration
Expose /healthz and /ready endpoints separately
