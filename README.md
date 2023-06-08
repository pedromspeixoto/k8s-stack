# k8s stack

Table of Contents
=================

   * [k8s stack](#k8s-stack)
      * [Overview](#overview)
      * [Prerequisites](#prerequisites)
      * [Setup](#setup)
         * [Manual](#manual)
         * [Script](#script)
      * [Demo](#demo)
      * [References](#references)

## Overview

This project has a simple demo on how to setup a serverless kubernetes cluster with a GitOps approach with the following tools:

- Kind (Kubernetes in Docker) - https://kind.sigs.k8s.io/docs/user/quick-start/
- Istio (Service Mesh) - https://istio.io/docs/setup/kubernetes/quick-start/
- Knative (Serverless) - https://knative.dev/docs/install/knative-with-kind/
- ArgoCD (GitOps) - https://argoproj.github.io/argo-cd/getting_started/

The CI pipeline builds the project and deploys the application to the cluster using ArgoCD. The application will be deployed to a single k8s cluster (for demo purposes) to a different namespace depending on the triggering branch.

- The main branch will deploy to the `prod` namespace
- The develop branch will deploy to the `staging` namespace

A short diagram on how the CI/CD project will be structured is presented below (for the sake of the demo, the infrastructure deployment manifests are in the same repo as the application code but in a real world scenario they would be in a separate repo):

![CI/CD Pipeline Diagram](/assets/CICD_Pipeline_Diagram.png)


## Prerequisites

- Docker (https://www.docker.com/products/docker-desktop)
- Kind (https://kind.sigs.k8s.io/docs/user/quick-start/)
- Kubectl (https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- Istioctl (https://istio.io/latest/docs/setup/getting-started/)
- Helm (https://helm.sh/docs/intro/install/)

At the moment of this setup this was tested with the following versions:

- Docker: 20.10.14
- Kind: v0.19.0 go1.20.4 darwin/arm64 and kindest/node:v1.27.1
- Kubectl: 1.21.9
- Istioctl: 1.17.2
- Helm: 3.10.2

## Setup

### Manual

Please follow the below steps to start the kind local cluster with Istio and Knative installed.

1. Create the kind cluster

```bash
kind create cluster --name k8s-demo --config "./configs/kind-cluster.yaml"
```

2. Install Istio

```bash
istioctl install --set profile=demo -y
```

3. Install Knative

```bash
kubectl apply --filename
```

4. Install ArgoCD using Helm

Before we install the ArgoCD chart, we need to generate a Chart.lock file for Argo. We do this so that our dependency (the original argo-cd chart) can be rebuilt. This is important later when we let Argo CD manage this chart to avoid errors. We can do this by running the following commands:

```bash
helm repo add argo https://argoproj.github.io/argo-helm
helm dep update argocd/helm
```

We also need to create the namespace where we will install ArgoCD:

```bash
kubectl create namespace argocd
```

Now we can install ArgoCD with the following command:

```bash
helm install argo-cd argocd/helm --namespace argocd
```

Once Argo is installed you can forward the serve port with the following command:

```bash
kubectl port-forward svc/argo-cd-argocd-server -n argocd 8080:443
```

You can now access the ArgoCD UI with the following URL: https://localhost:8080

*Notes:*

The default username for ArgoCD is `admin`. The password is auto-generated and we can get it with the following command:

```bash
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
```

5. Finally, we can configure ArgoCD to manage the cluster for staging and prod namespaces. We can do this by running the following command:

```bash
helm template argocd/apps/staging/ | kubectl apply -f -
helm template argocd/apps/prod/ | kubectl apply -f -
```

### Script

- Alternatively you can use any of the scripts located under the `scripts` folder to setup the cluster, setup ArgoCD and delete the cluster.

Example to setup the cluster:

```bash
./scripts/1-setup-cluster.sh
```

## Demo

TBD

## References

- https://www.danielstechblog.io/
- https://medium.com/@s4l1h