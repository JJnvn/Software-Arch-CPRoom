# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app/backend

COPY . .

WORKDIR /app/backend/services/booking

RUN go mod download
RUN go build -o /booking-service ./cmd/server.go

# Run stage
FROM debian:bookworm-slim

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /booking-service .

# Expose gRPC port
EXPOSE 50051

# Run command
CMD ["./booking-service"]
