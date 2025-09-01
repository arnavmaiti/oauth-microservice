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
You can read the full list of changes [here]().

#### :rocket: New Features
* TODO

#### :bug: Bug Fixes
* Nothing here

## :wrench: Getting Started

### Build and Run With Docker
* `docker build -t oauth-microservice:latest .`
* `docker run -p 8080:8080 oauth-microservice:latest`

### Helm Chart
* For first install: `helm install auth-service ./charts/auth-service`
* `kubectl get pods`
* `kubectl port-forward svc/auth-service 8080:8080`
* For further updates: `helm upgrade --install auth-service ./charts/auth-service`

