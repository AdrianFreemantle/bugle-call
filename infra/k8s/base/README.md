# k8s/base

This folder contains reusable, GitOps-managed building blocks for the Bugle Call platform.

Each subdirectory defines a deployable unit (e.g., NATS, PostgreSQL, incident-api) using either raw Kubernetes manifests or HelmRelease YAMLs.

These components are referenced by environment overlays (e.g. `overlays/dev`, `overlays/staging`) and are expected to be managed by ArgoCD or `kubectl apply -k`.

Do not place CRDs or one-time setup scripts here.
