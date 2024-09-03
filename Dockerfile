# syntax = docker/dockerfile:1

FROM jdxcode/mise:latest AS build

ARG VERSION=development
ARG COMMIT=development

ENV VERSION=${VERSION}
ENV COMMIT_SHA=${COMMIT}

WORKDIR /app

# Install unzip dependency for bun
RUN apt-get update && apt-get install -y unzip
# Install Taskfile
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/bin

# Install language runtimes from .mise.toml
COPY /.mise.toml ./.mise.toml
RUN mise install && mise activate --shims bash

# Cache build dependencies
ENV GOCACHE=/root/.cache/go-build

# Cache Go modules
COPY core/go.mod core/go.sum ./core/
COPY package.json bun.lockb ./
COPY dashboard/package.json ./dashboard/
COPY tracker/package.json ./tracker/

RUN --mount=type=cache,target=${GOCACHE} \
	--mount=type=cache,target=/go/pkg/mod \
	cd core && go mod download && cd ..

RUN bun install --frozen-lockfile

# Copy the rest of the source code
COPY . .
RUN --mount=type=cache,target=${GOCACHE} ~/bin/task core:release:docker

# Build the final image
FROM gcr.io/distroless/cc-debian12

LABEL org.opencontainers.image.source=https://github.com/medama-io/medama \
	org.opencontainers.image.description="Cookie-free, privacy-focused website analytics." \
	org.opencontainers.image.licenses=Apache-2.0

ENV PORT=8080 \
	ANALYTICS_DATABASE_HOST=/app/data/me_analytics.db \
	APP_DATABASE_HOST=/app/data/me_app.db

WORKDIR /app

# Copy the binary
COPY --from=build /app/core/bin/main /app/bin/main

EXPOSE ${PORT}
CMD ["/app/bin/main", "start", "-env"]

