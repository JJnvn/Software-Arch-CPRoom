# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app/backend

COPY . .

WORKDIR /app/backend/services/approval

RUN go build -o /approval-service ./cmd/server.go

# Runtime stage
FROM debian:bookworm-slim

WORKDIR /root/

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /approval-service .

EXPOSE 8085 50052

CMD ["./approval-service"]
