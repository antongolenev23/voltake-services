FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth-service ./services/auth/cmd/auth/main.go


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/auth-service .
CMD ["./auth-service"]