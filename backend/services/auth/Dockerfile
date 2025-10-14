# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app/backend

COPY . .

WORKDIR /app/backend/services/auth

RUN go mod download
RUN go build -o /auth-service ./cmd/server.go

# Run stage with small image
FROM debian:bookworm-slim

WORKDIR /root/

# Install CA certificates so Go can verify HTTPS
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy binary from builder
COPY --from=builder /auth-service .

# Expose port
EXPOSE 8081

# Command to run
CMD ["./auth-service"]
