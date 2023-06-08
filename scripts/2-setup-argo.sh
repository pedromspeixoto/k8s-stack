#!/bin/bash

ARGO_NAMESPACE="argocd"
RELEASE_NAME="argo-cd"

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

# Check if the Argo CD Helm chart is already installed
helm list -n "$ARGO_NAMESPACE" | grep -q "$RELEASE_NAME"
if [ $? -eq 0 ]; then
  echo "Argo CD Helm chart is already installed in namespace $ARGO_NAMESPACE."
else
  echo "Installing Argo CD Helm chart in namespace $ARGO_NAMESPACE..."
  helm install "$RELEASE_NAME" argocd/helm --namespace "$ARGO_NAMESPACE"
fi

# Wait for ArgoCD to be ready
SERVICE_NAME="argo-cd-argocd-server"
TIMEOUT=300  # 5 minutes timeout (in seconds)

# Calculate the end time for the timeout
END_TIME=$((SECONDS + TIMEOUT))

# Wait for the service to be available
echo "Waiting for service $SERVICE_NAME in namespace $ARGO_NAMESPACE to be available..."

while [[ $SECONDS -lt $END_TIME ]]; do
  SERVICE_STATUS=$(kubectl get service "$SERVICE_NAME" -n "$ARGO_NAMESPACE" -o jsonpath='{.spec.clusterIP}')

  if [[ "$SERVICE_STATUS" != "<pending>" ]]; then
    echo "Service $SERVICE_NAME is now available. You can start port forwarding."
    break
  fi

  echo "Service $SERVICE_NAME is not available yet. Current status: $SERVICE_STATUS"
  sleep 5
done

if [[ $SECONDS -ge $END_TIME ]]; then
  echo "Timeout: Service $SERVICE_NAME did not become available within the specified time."
fi