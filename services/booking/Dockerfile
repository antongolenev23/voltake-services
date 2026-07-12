FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o booking-service ./services/booking/cmd/booking/main.go


FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache curl ca-certificates && update-ca-certificates

COPY --from=builder /app/booking-service .

CMD ["./booking-service"]