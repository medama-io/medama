# syntax = docker/dockerfile:1

# Setup Go builder
FROM golang:1.21 AS builder

# Install Taskfile
RUN go install github.com/go-task/task/v3/cmd/task@latest

# Copy go.mod and go.sum and download dependencies
WORKDIR /app
COPY go.mod go.sum ./

# Copy the rest of the source code and build
COPY . .
RUN task build

# Setup lightweight container
FROM alpine:latest

# Copy the binary from the builder stage
COPY --from=builder /app/bin/main /app/main

# Run the binary
EXPOSE 8080
CMD ["/app/main", "start", "-debug"]

