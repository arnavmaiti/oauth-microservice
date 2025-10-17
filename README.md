# OAuth Microservice
![GitHub go.mod Go version (branch)](https://img.shields.io/github/go-mod/go-version/arnavmaiti/oauth-microservice)
![GitHub Release](https://img.shields.io/github/v/release/arnavmaiti/oauth-microservice)


## Table of Contents
* [Introduction](#bulb-introduction)
* [Built With](#package-built-with)
* [What's New](#sparkles-whats-new)
* [Getting Started](#wrench-getting-started)

## :package: Built With
* [Go](https://go.dev/)

## :sparkles: What's New

### Version 3.0.0 (Latest)
You can read the full list of changes [here](https://github.com/arnavmaiti/oauth-microservice/wiki/Version-3.0.0).

#### :rocket: New Features
* Split auth and user services to respective modules
* Inter-pod communication
* All OAuth APIs now have proper functionalities and request and response structure. The details can be found [here](https://github.com/arnavmaiti/oauth-microservice/wiki/OAuth-Endpoints)

#### :bug: Bug Fixes
* Updated all APIs to conform to [RFC 6749](https://datatracker.ietf.org/doc/html/rfc6749)

## :bulb: Introduction
Go-based OAuth 2.0 + OIDC Authorization Server with a React frontend, designed to run as microservices on Kubernetes
## :rocket: Release Notes
* [v3.0.0](https://github.com/arnavmaiti/oauth-microservice/wiki/Version-3.0.0)
* [v2.0.0](https://github.com/arnavmaiti/oauth-microservice/wiki/Version-2.0.0)
* [v1.0.0](https://github.com/arnavmaiti/oauth-microservice/wiki/Version-1.0.0)

## :wrench: Getting Started

### Build and Run With Docker Compose
* Docker compose is used to build the development services
* Docker images can be built using `docker-compose build`. This builds Auth service and User service.
* Container orchestration can then be turned on by `docker-compose up`.
  * Auth service is available on port 8080
  * User service is available on port 8081
  * PostGRES service is available on post 5432
* Subsequent builds and start can be done by `docker-compose up --build`.
* In order to rebuild latest pod delete an existing pod by using `kubectl delete pod <pod-name>`
* In order to rebuild from scratch delete the complete namespace `kubectl delete namespace default`
* Access PostGRES container by `docker exec -it postgres psql -U authuser -d authdb`

### Helm Chart

#### General
* To check the pods: `kubectl get pods`
* Access PostGRES database using: `kubectl exec -it <postgres_pod_name> -- psql -U authuser -d authdb`

#### Complete Deployment
* Build the containers locally `docker-compose build`
* Update dependencies `helm dependency update charts`
* For first run `helm install oauth-microservice charts`
* For further updates: `helm upgrade --install oauth-microservice charts`

#### Auth Service
* For first install: `helm install auth-service charts/auth-service`
* To start the server at http://localhost:8080: `kubectl port-forward svc/auth-service 8080:8080`
  * View health API at `GET http://localhost:8080/health`
  * View ready API at `GET http://localhost:8080/health`
* For further updates: `helm upgrade --install auth-service ./charts/auth-service`

#### User Service
* For first install: `helm install user-service charts/user-service`
* To start the server at http://localhost:8081: `kubectl port-forward svc/user-service 8081:8081`
  * View users API at `GET http://localhost:8081/users`
* For further updates: `helm upgrade --install user-service ./charts/user-service`

#### Migrations
* For first install: `helm install migrations charts/migrations`
* For further updates: `helm upgrade --install migrations ./charts/migrations`

### (Optional) Ingress on local
* Install ingress-nginx
```
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm uninstall ingress-nginx -n ingress-nginx
helm install ingress-nginx ingress-nginx/ingress-nginx --namespace ingress-nginx --create-namespace --set controller.admissionWebhooks.enabled=false
```
* Setup hosts file with the following line `127.0.0.1 auth.local`
* Upgrade helm
* Health and readiness APIs should be available
```
https://auth.local/health
https://auth.local/ready
```

### (Temporary till client register is implemented) How to create a client in PostGRES
* Get the PostGRES pod by `kubectl get pods`
* Execute bash command `kubectl exec -it <postgres-0> -- bash`
* Run psql command line `psql -U authuser -d authdb`
* Run the following SQL command to create a temporary client
```
INSERT INTO oauth_clients (id, client_id, client_secret, redirect_uris, scopes, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    'client123',
    'secret123',
    ARRAY['http://localhost:8080/callback'],
    ARRAY['openid'],
    NOW(),
    NOW()
);
```

### How to create and test flow
[Read it here](https://github.com/arnavmaiti/oauth-microservice/wiki/Version-2.0.0#book-how-to-create-user-and-test-flow)


