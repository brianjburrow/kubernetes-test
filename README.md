# Get set up

## Start minikube and add the minikube registry

```bash
minikube delete -h
minikube start --driver=docker
minikube status
minikube addons enable registry
minikube addons enable dashboard
minikube addons enable 
kubectl get svc -n kube-system registry # Copy the IP into the ``deploy.sh`` file and ./k8s/*.yaml deployment files.
```

Note down the Cluster IP of the registry, and copy it into

- the `deploy.sh` file
- ./k8s/*.yaml deployment files

Then, you'll need to open two new terminals and port-forward the frontend and backend ports.

```bash
kubectl port-forward deployment/frontend-deployment 8080:80
kubectl port-forward deployment/backend-deployment 8081:80
```
