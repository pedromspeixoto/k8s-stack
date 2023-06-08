#!/bin/bash
set -e

ARGO_NAMESPACE="argocd"

helm template argocd/apps/staging/ --namespace $ARGO_NAMESPACE | kubectl apply -f -
helm template argocd/apps/prod/ --namespace $ARGO_NAMESPACE | kubectl apply -f -