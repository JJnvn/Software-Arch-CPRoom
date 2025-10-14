# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app/backend

COPY . .

WORKDIR /app/backend/services/room

RUN go mod download
RUN go build -o /room-service ./cmd/server.go

# Run stage
FROM debian:bookworm-slim

WORKDIR /root/
COPY --from=builder /room-service .

EXPOSE 8082

CMD ["./room-service"]
