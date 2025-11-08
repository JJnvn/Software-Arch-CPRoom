# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app/backend

# Copy the entire backend workspace so local replaces resolve
COPY . .

WORKDIR /app/backend/services/auth

RUN go build -o /auth-service ./cmd/server.go

# Runtime stage
FROM debian:bookworm-slim

WORKDIR /root/

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /auth-service .

EXPOSE 8081

CMD ["./auth-service"]
