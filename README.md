# Get set up

## Start minikube and add the minikube registry

```bash
minikube delete -h
minikube start --driver=docker
minikube status
minikube addons enable registry
kubectl get svc -n kube-system registry
```

Note down the Cluster IP of the registry

## Start pushing the images to the minikube registry

Point the docker CLI to the minikube docker daemon

```bash
eval $(minikube -p minikube docker-env)
```

Build the images

```bash
docker build -t 10.110.159.84:80/client-image:latest ./client 
docker build -t 10.110.159.84:80/backend-image:latest ./backend-go 
docker build -t 10.110.159.84:80/mysql-image:latest ./mysql
docker build -t 10.110.159.84:80/frontend-image:latest ./frontend
```

Push the images to the registry

```bash
docker push 10.110.159.84:80/backend-image:latest
docker push 10.110.159.84:80/client-image:latest
docker push 10.110.159.84:80/mysql-image:latest
docker push 10.110.159.84:80/frontend-image:latest
```

Before applying the manifests, make sure all the images are present

```bash
curl 10.110.159.84:80/v2/_catalog
```

Apply the kubernetes manifests

```bash
kubectl apply -f k8s/
```

Verify the status

```bash
kubectl get pods
```

```bash
kubectl get svc
```

If you can't connect to the front end, there may be issues, try port forwarding to the container

```bash
kubectl port-forward pod/frontend-deployment-6bfccdd6f7-gjbjh 8080:80
```

And then connect via local host

```bash
http://localhost:8080/
```

You can roll out restarts like this

```bash
kubectl rollout restart deployment frontend-deployment
```

Which updates the pod name
