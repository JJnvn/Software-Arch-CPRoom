# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of the code
COPY . .

# Build binary
RUN go build -o user-service ./cmd/server.go

# Run stage
FROM debian:bookworm-slim

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/user-service .

# Expose gRPC port
EXPOSE 50051

# Run command
CMD ["./user-service"]
