#!/bin/bash

#if my-cluster exists, delete it
if kind get clusters | grep -q bluepill; then
  kind delete cluster --name bluepill
fi

#build Docker image
docker build -t docker.io/hoeg/bluepill:latest .

# Create a Kind cluster
kind create cluster --name bluepill

# Load the Docker image into the Kind cluster
kind load docker-image docker.io/hoeg/bluepill:latest --name bluepill

# apply all the resources in this folder except test-deployment.yaml
kubectl apply -f secret.yaml
kubectl apply -f deploy/whitelist-config.yaml
kubectl apply -f deploy/configuration.yaml
kubectl apply -f deploy/deployment.yaml
kubectl apply -f deploy/service.yaml
#kubectl apply -f deploy/admission-webhook.yaml

# wait for the pods to be ready
kubectl wait --for=condition=ready --timeout=5m pod --all