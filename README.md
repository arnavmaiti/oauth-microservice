# OAuth Microservice
![GitHub go.mod Go version (branch)](https://img.shields.io/github/go-mod/go-version/arnavmaiti/oauth-microservice/v1.0.0)

## Table of Contents
* [Introduction](#bulb-introduction)
* [Built With](#package-built-with)
* [What's New](#sparkles-whats-new)
* [Getting Started](#wrench-getting-started)

## :bulb: Introduction
Go-based OAuth 2.0 + OIDC Authorization Server with a React frontend, designed to run as microservices on Kubernetes

## :package: Built With
* [Go](https://go.dev/)

## :sparkles: What's New

### Version 1.0.0 (Latest)
You can read the full list of changes [here](https://github.com/arnavmaiti/oauth-microservice/wiki/Version-1.0.0).

#### :rocket: New Features
* TODO

#### :bug: Bug Fixes
* Nothing here

## :wrench: Getting Started

### Build and Run With Docker
* We will need to establish a network for service and a sample PostGRES to work together
* `docker network create mynet`
* Run PostGRES sample container `docker run -d --name pg --network mynet -e POSTGRES_PASSWORD=changeme123 postgres:15`
* Build latest docker `docker build -t oauth-microservice:latest .`
* Run the container `docker run --network mynet -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=changeme123 -e POSTGRES_HOST=pg -e POSTGRES_PORT=5432 -e POSTGRES_DB=postgres -p 8080:8080 oauth-microservice:latest`
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

