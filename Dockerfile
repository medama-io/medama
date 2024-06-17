# syntax = docker/dockerfile:1

FROM jdxcode/mise:latest AS build

# Install unzip dependency for bun
RUN apt-get update && apt-get install -y unzip

# Install runtimes - Temporarily include node for google-closure-compiler for tracker
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

LABEL org.opencontainers.image.source=https://github.com/medama-io/medama
LABEL org.opencontainers.image.description="Cookie-free, privacy-focused website analytics."
LABEL org.opencontainers.image.licenses=Apache-2.0

ENV PORT=8080
ENV ANALYTICS_DATABASE_HOST=/app/data/me_analytics.db
ENV APP_DATABASE_HOST=/app/data/me_app.db

WORKDIR /app

# Copy the binary
COPY --from=build /app/core/bin/main /app/bin/main

EXPOSE ${PORT}
CMD ["/app/bin/main", "start", "-env"]

