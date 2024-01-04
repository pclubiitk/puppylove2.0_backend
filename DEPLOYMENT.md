Well be using minkube to setup our local k8 cluster. Install minikube[https://minikube.sigs.k8s.io/docs/start/] and kubectl[https://kubernetes.io/docs/tasks/tools/].


```
minikube start
minkube dashboard
```
Now switch to another terminal
## Running the database pods:


First we will get the database running with persistant volume:
```
kubectl apply -f kubernetes/deployments/db-core-deployment.yaml

```
Exposing the postgres container:
```
kubectl apply -f kubernetes/services/db-core-service.yaml
kubectl get svc db-core
```

We will now create/mount the volumes:
```
kubectl apply -f kubernetes/persistent-volumes/db-core-volume.yaml
```
We can check if the db pods are running:
```
kubectl get pods
```
You can enter into the postgres shell by executing the following command:
```
kubectl exec -it db-core-deployment-57d8f9746c-5dfnm psql "postgresql://postgres:postgres@localhost:5432/puppylove"

```

## Running the backend pods:
First we will apply the environment configuration stored in ` kubernetes/configs/backend-core-config.yaml` which is of the following format(you can change this):

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: backend-core-config
data:
  POSTGRES_HOST : "db-core"
  POSTGRES_PORT : "5432"
  POSTGRES_PASSWORD : "postgres"
  POSTGRES_DB : "puppylove"
  POSTGRES_USER : "postgres"
  CFG_ADMIN_PASS : "something"

  ADMIN_ID : "admin1"
  ADMIN_PASS : "admin2"

  USER_JWT_SIGNING_KEY : "something"
  HEART_JWT_SIGNING_KEY : "something2"
  DOMAIN : "localhost"

  EMAIL_ID : "hello@iitk.ac.in"
  EMAIL_PASS : "hello"
```
```
kubectl apply -f kubernetes/configs/backend-core-config.yaml 
```
Now deploy the backend pod with the env variables above:
```
kubectl apply -f kubernetes/deployments/backend-core-deployments.yaml 
```

We can also expose this backend to the outside network:

```
kubectl apply -f kubernetes/services/backend-service.yaml
kubectl get svc backend-core
```
You can enter into the container shell:
```
kubectl exec -it backend-core-deployment-58cf59cfdd-c9xcb /bin/bash

```

Now you should have a running backend and database. Enter into the backend shell and try to execute some api's using curl.

### References:

https://medium.com/@wu.victor.95/deploying-with-kubernetes-d3a9e9aad767
https://medium.com/@wu.victor.95/deploying-with-kubernetes-pt-2-645434c84af3