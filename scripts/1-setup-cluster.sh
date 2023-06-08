#!/bin/bash
set -e

## Creating a kind cluster
echo "Create a kind cluster..."
kind create cluster --name k8s-demo --config ./configs/kind-cluster.yaml

## Set the kubectl context to the created kind cluster
echo "Set the kubectl context to the created kind cluster..."
kubectl cluster-info --context kind-k8s-demo

## Installing Istio on the created kind cluster
echo "Installing Istio on the created kind cluster..."
istioctl install -f ./configs/istio.yaml -y

## Installing required custom resources for knative serving component on the created kind cluster
echo "Installing knative on the created kind cluster..."
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.10.1/serving-crds.yaml

## Installing the core components of knative serving on the created kind cluster
echo "Installing the core components of knative serving on the created kind cluster..."
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.10.1/serving-core.yaml

## Installing the knative Istio controller on the created kind cluster
echo "Installing the knative Istio controller on the created kind cluster..."
kubectl apply -f https://github.com/knative/net-istio/releases/download/knative-v1.10.0/net-istio.yaml

## Enabling Istio automatic sidecar injection for the knative-serving namespace
echo "Enabling Istio automatic sidecar injection for the knative-serving namespace..."

## Configuring DNS to use sslip.io - sslip.io provides a wildcard DNS setup that will automatically resolve to the IP address you put in front of sslip.io.
##kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.10.1/serving-default-domain.yaml
kubectl patch configmap/config-domain --namespace knative-serving --type merge --patch '{"data":{"127.0.0.1.nip.io":""}}'

## Configure mTLS for knative serving which secures service-to-service communication within the cluster
echo "Configure mTLS for knative serving which secures service-to-service communication within the cluster..."
kubectl label namespace knative-serving istio-injection=enabled
kubectl apply -f configs/knative-mtls-config.yaml

## Create namespaces
echo "Create default namespaces for different environments..."
kubectl create namespace staging
kubectl create namespace prod
kubectl label namespace staging istio-injection=enabled
kubectl label namespace prod istio-injection=enabled