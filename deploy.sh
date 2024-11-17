#!/bin/bash

# Set image tag (use a timestamp or a unique identifier for each build)
IMAGE_TAG="latest"
REGISTRY_IP=$(minikube ip)

# Point the shell to minikube's docker-daemon by setting env variables
eval $(minikube -p minikube docker-env)

# Build all images
docker build -t 10.105.247.24:80/client-image:$IMAGE_TAG ./client
docker build -t 10.105.247.24:80/backend-image:$IMAGE_TAG ./backend-go
docker build -t 10.105.247.24:80/frontend-image:$IMAGE_TAG ./frontend
# docker build -t client-image:$IMAGE_TAG ./client
# docker build -t backend-image:$IMAGE_TAG ./backend-go
# docker build -t frontend-image:$IMAGE_TAG ./frontend

# Push images to the registry
docker push 10.105.247.24:80/client-image:$IMAGE_TAG
docker push 10.105.247.24:80/backend-image:$IMAGE_TAG
docker push 10.105.247.24:80/frontend-image:$IMAGE_TAG

# Update the Kubernetes Deployment for each service with the new image
kubectl apply -f k8s/

# Restart the deployments
kubectl rollout restart deployment client-deployment
kubectl rollout restart deployment backend-deployment
kubectl rollout restart deployment mysql
kubectl rollout restart deployment frontend-deployment

# Wait for the frontend deployment to be fully rolled out
kubectl rollout status deployment client-deployment
kubectl rollout status deployment backend-deployment
kubectl rollout status deployment mysql
kubectl rollout status deployment frontend-deployment

# Port-forward to the frontend pod (adjust port as necessary)
kubectl port-forward deployment/frontend-deployment 8080:80
