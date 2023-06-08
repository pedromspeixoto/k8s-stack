#!/bin/bash
set -e

ARGO_NAMESPACE="argocd"

# Create Chart.lock file for ArgoCD - so that our dependency (the original argo-cd chart) can be rebuilt
helm repo add argo https://argoproj.github.io/argo-helm
helm dep update argocd/helm

# Check if the namespace exists
if kubectl get namespace "$ARGO_NAMESPACE" >/dev/null 2>&1; then
  echo "Namespace $ARGO_NAMESPACE already exists. Skipping creation."
else
  echo "Creating namespace $ARGO_NAMESPACE..."
  kubectl create namespace "$ARGO_NAMESPACE"
fi

# Install argocd using Helm
helm install argo-cd argocd/helm --namespace $ARGO_NAMESPACE