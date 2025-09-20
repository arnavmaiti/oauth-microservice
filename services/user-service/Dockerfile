# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /application

# Install git (needed for go mod)
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o user-service ./services/user-service

# Runtime stage
FROM alpine:latest
WORKDIR /application

RUN apk add --no-cache ca-certificates libc6-compat

COPY --from=builder /application/user-service .

EXPOSE 8081

CMD ["./user-service"]
