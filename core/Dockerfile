# syntax = docker/dockerfile:1

# Setup Go builder
FROM golang:1.21 AS backend

# Install Taskfile
RUN go install github.com/go-task/task/v3/cmd/task@latest

COPY . /core
WORKDIR /core

# Build the core
RUN go mod download
RUN task build

EXPOSE 8080

CMD ["/core/bin/main", "start", "-debug"]

