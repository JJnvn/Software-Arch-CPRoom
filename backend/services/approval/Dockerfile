# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app/backend

COPY . .

WORKDIR /app/backend/services/approval

RUN go mod tidy
RUN go build -o /approval-service ./cmd/server.go

# Run stage
FROM debian:bookworm-slim

WORKDIR /root/

COPY --from=builder /approval-service .

EXPOSE 8085

CMD ["./approval-service"]
