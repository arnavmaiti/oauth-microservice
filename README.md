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

### Version 2.0.0 (Latest)
You can read the full list of changes [here](https://github.com/arnavmaiti/oauth-microservice/wiki/Version-2.0.0).

#### :rocket: New Features
* Endpoints implemented:
```
Endpoint           | Method | Purpose
/register          | POST   | Create new users in Postgres
/authorize         | GET    | Initiate OAuth2 authorization code flow
/token             | POST   | Exchange authorization code for access token & refresh token; supports refresh token flow
/introspect        | POST   | Validate access token and return metadata (user, scopes, expiry)
/revoke (optional) | POST   | Revoke access or refresh tokens
```

#### :bug: Bug Fixes
* Nothing here

## :bulb: Introduction
Go-based OAuth 2.0 + OIDC Authorization Server with a React frontend, designed to run as microservices on Kubernetes
## :rocket: Release Notes
* [v1.0.0](https://github.com/arnavmaiti/oauth-microservice/wiki/Version-1.0.0)
* [v2.0.0](https://github.com/arnavmaiti/oauth-microservice/wiki/Version-2.0.0)

## :wrench: Getting Started

### Build and Run With Docker
* We will need to establish a network for service and a sample PostGRES to work together
* `docker network create mynet`
* Run PostGRES sample container `docker run -d --name pg --network mynet -e POSTGRES_PASSWORD=changeme123 postgres:15`
* Build latest docker `docker build -t oauth-microservice:latest .`
* Run the container `docker run --network mynet -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=changeme123 -e POSTGRES_HOST=pg -e POSTGRES_PORT=5432 -e POSTGRES_DB=postgres -p 8080:8080 oauth-microservice:latest`
* Please note, in order to get the latest local pod in helm, use `kubectl delete pod <oauth-microservice-pod>`
* You should see 
```
OAuth server running on :8080
2025/09/05 15:32:41 Successfully connected to database
```

### Helm Chart
* For first install: `helm install auth-service ./charts/auth-service`
* To check the pods: `kubectl get pods`
* To start the server at http://localhost:8080: `kubectl port-forward svc/oauth-microservice 8080:8080`
> * View health API at `GET http://localhost:8080/health`
* For further updates: `helm upgrade --install auth-service ./charts/auth-service`
* For updates with migrations: `helm upgrade --install auth-service ./charts/auth-service --set migrations.enabled=true` or `helm upgrade --install migrations ./charts/auth-service/charts/migrations`
* Access PostGRES database using: `kubectl exec -it <postgres_pod_name> -- psql -U authuser -d authdb`

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


