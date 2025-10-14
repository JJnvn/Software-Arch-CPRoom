FROM golang:1.24-alpine AS builder

WORKDIR /app/backend

COPY . .

WORKDIR /app/backend/services/notification

RUN go mod tidy && go build -o /notification-service ./cmd/server.go

FROM debian:bookworm-slim

WORKDIR /root/

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /notification-service .

EXPOSE 8084

CMD ["./notification-service"]
