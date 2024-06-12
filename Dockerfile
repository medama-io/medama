# syntax = docker/dockerfile:1

FROM jdxcode/mise:latest AS build

# Install unzip dependency for bun
RUN apt-get update && apt-get install -y unzip

# Install runtimes
RUN mise use -g node@20
RUN mise use -g bun@latest
RUN mise use -g go@1.22

# Install Taskfile
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/bin

# Cache build dependencies
ENV GOCACHE=/root/.cache/go-build

# Cache Go modules
WORKDIR /app
COPY /core/go.mod /core/go.sum ./core/

WORKDIR /app/core
RUN --mount=type=cache,target=/root/.cache/go-build \
	--mount=type=cache,target=/go/pkg/mod \
	go mod download

# Go back to root directory and install JavaScript dependencies
WORKDIR /app
COPY package.json bun.lockb ./
COPY dashboard/package.json ./dashboard/
COPY tracker/package.json ./tracker/
RUN bun install --frozen-lockfile

# Copy the rest of the source code
COPY . .
RUN --mount=type=cache,target="/root/.cache/go-build" ~/bin/task core:release

# Build the final image
FROM gcr.io/distroless/cc-debian12

WORKDIR /app

# Copy the binary
COPY --from=build /app/core/bin/main /app/bin/main

# Run the binary
EXPOSE 8080
CMD ["/app/bin/main", "start", "-level=debug"]

