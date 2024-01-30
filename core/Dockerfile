# syntax = docker/dockerfile:1

# Setup Go builder
FROM golang:1.21

# Cache build dependencies
ENV GOCACHE=/root/.cache/go-build

# Install Taskfile
RUN --mount=type=cache,target=/root/.cache/go-build \
	go install github.com/go-task/task/v3/cmd/task@latest

# Copy go.mod and go.sum and download dependencies
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build \
	go mod download

# Copy the rest of the source code and build
COPY . .
RUN --mount=type=cache,target="/root/.cache/go-build" task build

# Run the binary
EXPOSE 8080
CMD ["/app/bin/main", "start", "-debug"]

