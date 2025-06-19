# k8s/bootstrap

This folder contains one-time bootstrap resources for initializing the cluster before GitOps takes over.

It includes:
- CRDs required by operators (e.g., External Secrets)
- Cluster-wide resources like `ClusterSecretStore`
- Initial namespace and RBAC setup
- GitOps bootstrap YAMLs (e.g., ArgoCD HelmRepository)

These should be applied manually or via an install script, not managed by ArgoCD.
