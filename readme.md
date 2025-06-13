# Bugle Call

**Bugle Call is a full-stack incident intelligence platform built for testing architecture patterns, GitOps delivery, and platform resilience in real-world conditions.**

## Purpose

This project is designed to simulate the realities of cloud-native architecture in a controlled, measurable environment. It focuses on the kind of work architects often design but rarely get to implement fully themselves: end-to-end system behavior under load, cost-awareness, secure access control, observability, and practical resilience.

The goal is to deepen architectural skill through hands-on implementation, using modern tools and patterns across the full software lifecycle.

## What We're Building

Bugle Call is a distributed system that accepts operational incidents, enriches and classifies them using machine learning, and provides secure visibility into trends, anomalies, and system behavior. It integrates intelligent processing, real-time observability, and infrastructure automation.

It includes AI capabilities for incident categorization and cost anomaly detection, using lightweight local models that enhance system insight without overcomplicating the architecture.

## Core Components
[View the full architecture](docs/architecture.md)

| Component                    | Description                                                             |
|------------------------------|-------------------------------------------------------------------------|
| Incident API (.NET)          | Accepts HTTP POSTs, validates incidents, emits events                   |
| Async Processor (Go)         | Processes messages, triggers enrichment workflows                       |
| Enrichment Service (Python)  | Tags, classifies, and prioritizes incidents using lightweight ML        |
| Cost Anomaly Detector        | Analyzes metrics from Prometheus and Kubecost to flag unusual behavior  |
| PostgreSQL                   | Central structured incident storage                                     |
| Admin UI (React)             | Internal dashboard for browsing, filtering, and analyzing incidents     |
| Auth Service (Node)          | Handles user sessions and JWT-based access control                      |
| Load Generator               | Simulates realistic traffic and incident volumes                        |
| Monitoring Stack             | Prometheus, Grafana, Alertmanager                                       |
| GitOps Pipeline              | GitHub Actions and ArgoCD for delivery and promotion                    |
| Infrastructure               | Managed with Terraform for kind, AKS, and EKS environments              |
| GraphQL API (Node)           | Unified query layer for the UI                                          |
| Raw Payload Store (MongoDB)  | Stores raw incident JSON                                                | 

## AI Integrations

| Feature                    | Description                                                               |
|----------------------------|---------------------------------------------------------------------------|
| Incident Categorization     | Uses NLP to assign `category` and `severity` fields                       |
| Cost Anomaly Detection      | Detects unusual usage or cost patterns from Prometheus and Kubecost data  |

These integrations are fully self-hosted, resource-light, and designed to highlight practical AI enrichment within a platform context.

## Deployment Targets

| Environment        | Use Case                                                          |
|--------------------|--------------------------------------------------------------------|
| Local (`kind`, `k3d`) | Developer testing, pipeline validation                        |
| AKS (Azure)        | Enterprise IAM, regional cost modeling                             |
| EKS (AWS)          | Spot pricing, IRSA, multi-cloud comparison                         |

Environments are configured using Kustomize overlays and deployed via ArgoCD. Secrets and cloud-specific configurations are separated by layer.

## Platform Capabilities

### Simulated Load and Scaling
- Burst and idle traffic patterns
- Autoscaling and HPA configurations
- Fault injection for stress testing

### Multi-Environment GitOps Pipeline
- ArgoCD with GitHub Actions
- Kustomize overlays for dev, staging, and prod
- Promotion based on pull request workflow

### Secrets and IAM Strategy
- Azure Managed Identity and AWS IRSA support
- RBAC policies per namespace
- External Secrets and Vault integration planned

### Tunable Cost Scenarios
- Kubecost with environment-specific profiling
- Anomaly detection on resource usage
- Cross-cluster comparisons for pricing strategy

### Operational Resilience
- Liveness and readiness probes
- Message retry logic and visibility into failure modes
- Metrics, logs, and alerting built-in

## Technologies Used

| Category            | Stack                                                             |
|---------------------|-------------------------------------------------------------------|
| Infrastructure      | Kubernetes, ArgoCD, Kustomize, Terraform                          |
| Programming         | .NET 8, Go, Python 3.12, React , Add GraphQL, Apollo Server       |
| Messaging           | NATS or Redis Streams                                             |
| Monitoring          | Prometheus, Grafana, Alertmanager                                 |
| Cost Management     | Kubecost, anomaly detection models                                |
| IAM & Security      | JWT, RBAC, Azure Identity, AWS IRSA                               |
| Machine Learning    | spaCy, scikit-learn, pandas (local models only)                   |
| CI/CD               | GitHub Actions, Docker                                            |
| Persistence         | MongoDB, Postgres                                                 |

## Why This Project?

Architects often rely on others to implement their designs. This project reverses that. It puts the architect directly in the build path, exploring trade-offs, failures, and edge cases through real implementation.

Bugle Call is built to validate architectural patterns, expose system behavior under real conditions, and sharpen judgment through measurable feedback. It is a space to think clearly, build deliberately, and grow with intent.