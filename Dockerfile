# Stage 1: Build
FROM golang:1.25-alpine AS builder
WORKDIR /application

# Install dependencies
RUN apk add --no-cache git

# Cache Go modules
COPY go.mod go.sum* ./
RUN go mod download

# Copy code and build
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o oauth-microservice ./cmd/server

# Stage 2: Run
FROM alpine:latest
WORKDIR /application
COPY --from=builder /application/oauth-microservice .

EXPOSE 8080

# Run binary
CMD ["./oauth-microservice"]
