# Use the official Go 1.25-alpine image
FROM golang:1.25-alpine

WORKDIR /app

# Install bash (if necessary for your script)
RUN apk add --no-cache bash build-base

# Set up Go workspace
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

# Install air for live reloading
RUN go install github.com/air-verse/air@latest

COPY ./backend/go.mod ./backend/go.sum ./

RUN go mod tidy

