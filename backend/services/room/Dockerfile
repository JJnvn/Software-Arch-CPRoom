# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o room-service ./cmd/server.go

# Run stage
FROM debian:bookworm-slim

WORKDIR /root/
COPY --from=builder /app/room-service .

EXPOSE 8082

CMD ["./room-service"]
