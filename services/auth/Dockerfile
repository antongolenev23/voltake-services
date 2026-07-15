FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o auth-service ./services/auth/cmd/auth/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o seed-admin ./services/auth/cmd/seed-admin/main.go


FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache curl ca-certificates && update-ca-certificates

RUN set -eux; \
    curl -fL -o /usr/local/bin/grpc_health_probe \
    https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.24/grpc_health_probe-linux-amd64; \
    chmod +x /usr/local/bin/grpc_health_probe; \
    ls -la /usr/local/bin/grpc_health_probe

COPY --from=builder /app/auth-service .
COPY --from=builder /app/seed-admin .

CMD ["./auth-service"]