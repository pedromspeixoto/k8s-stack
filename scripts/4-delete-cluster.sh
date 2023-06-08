#!/bin/bash
set -e

## Cleaning kind cluster
echo "Deleting kind cluster..."
kind delete cluster --name k8s-demo