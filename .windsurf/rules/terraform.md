---
trigger: glob
globs: infra/terraform/**/*.tf
---

# File Layout
Each module must contain main.tf, variables.tf, outputs.tf
Use a clear path structure: /modules, /env/dev, /env/staging, etc.
Prefer module reuse over inline resource duplication

# Formatting and Linting
Run terraform fmt on every change
Use tflint to validate naming, security, and cloud-specific rules
Use terraform validate before committing

# Tagging and Metadata
All resources must have tags: environment, component, owner
Use locals block to define shared tags and apply them uniformly

# Secrets and IAM
Never define secret strings in .tf files
Use data sources to pull from secret stores (e.g., aws_secretsmanager_secret)
Prefer short-lived credentials via assumed roles or Managed Identity

# Backend and State
Use local backend only for kind/dev
All cloud environments must use remote state (e.g., S3 + DynamoDB, Azure Blob)
Do not commit .tfstate or .terraform directories
