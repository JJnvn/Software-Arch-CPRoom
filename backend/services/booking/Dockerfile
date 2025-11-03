# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app/backend

COPY . .

WORKDIR /app/backend/services/booking

RUN go build -o /booking-service ./cmd/server.go

# Run stage
FROM debian:bookworm-slim

WORKDIR /root/

COPY --from=builder /booking-service .

EXPOSE 8083 50051

CMD ["./booking-service"]
