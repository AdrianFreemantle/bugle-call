# k8s/overlays/dev

This overlay defines the full set of components and config for the `dev` environment.

It composes from reusable bases (in `../base`) and may apply environment-specific patches, configs, and Helm values overrides.

Apply with:
```bash
kubectl apply -k k8s/overlays/dev
```

Later, ArgoCD will continuously sync this overlay in a GitOps workflow.
