---
trigger: glob
globs: infra/k8s/**/*.yaml
---

# YAML Structure
Use 2-space indentation
Always use LF line endings and UTF-8 encoding
Separate multiple resources with `---`
Order fields consistently: apiVersion, kind, metadata, then spec

# Labels and Naming
Apply Kubernetes recommended labels: app.kubernetes.io/name, component, instance
Use lowercase kebab-case for resource names

# GitOps Structure
Follow Kustomize base/overlay pattern
All overlays must reference a shared base
Do not hardcode image tags or environment-specific values
Use patches or configMapGenerator where appropriate

# Secrets and ESO
Never commit native Kubernetes Secret objects
All secrets must be declared via ExternalSecret
Each namespace requiring secrets must define its own ExternalSecret
Secrets must sync from a centralized ClusterSecretStore

# Metrics and Scraping
All long-lived services must expose a /metrics endpoint
Add prometheus.io/scrape: "true" and prometheus.io/path annotations on Services or Pods
Avoid custom ports for metrics unless documented
