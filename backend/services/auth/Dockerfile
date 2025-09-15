# Use Go official image for building
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the binary
RUN go build -o auth-service ./cmd/server.go

# Run stage with small image
FROM debian:bookworm-slim

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/auth-service .

# Expose port
EXPOSE 8081

# Command to run
CMD ["./auth-service"]
