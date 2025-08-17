# syntax = docker/dockerfile:1@sha256:38387523653efa0039f8e1c89bb74a30504e76ee9f565e25c9a09841f9427b05

FROM debian:bookworm@sha256:731dd1380d6a8d170a695dbeb17fe0eade0e1c29f654cf0a3a07f372191c3f4b AS build

ARG VERSION=development
ARG COMMIT_SHA=development

ENV VERSION=${VERSION}
ENV COMMIT_SHA=${COMMIT_SHA}

RUN apt-get update  \
    && apt-get -y --no-install-recommends install  \
        curl git ca-certificates build-essential unzip \
    && rm -rf /var/lib/apt/lists/*

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
ENV MISE_DATA_DIR="/mise"
ENV MISE_CONFIG_DIR="/mise"
ENV MISE_CACHE_DIR="/mise/cache"
ENV MISE_INSTALL_PATH="/usr/local/bin/mise"
ENV PATH="/mise/shims:$PATH"

RUN curl https://mise.run | sh

WORKDIR /app

# Install language runtimes from .mise.toml
COPY mise.toml ./mise.toml
RUN mise trust -a -y && mise settings set experimental true && mise install && mise activate --shims bash

# Cache build dependencies
ENV GOCACHE=/root/.cache/go-build

# Cache Go modules
COPY core/go.mod core/go.sum ./core/
COPY package.json bun.lock ./
COPY dashboard/package.json ./dashboard/
COPY tracker/package.json ./tracker/

RUN --mount=type=cache,target=${GOCACHE} \
	--mount=type=cache,target=/go/pkg/mod \
	cd core && go mod download && cd ..

RUN bun install --frozen-lockfile

# Verify environment variables
RUN echo "VERSION=${VERSION}" && echo "COMMIT_SHA=${COMMIT_SHA}"

# Copy the rest of the source code
COPY . .
RUN --mount=type=cache,target=${GOCACHE} task core:release:docker

# Build the final image
FROM gcr.io/distroless/cc-debian12@sha256:00cc20b928afcc8296b72525fa68f39ab332f758c4f2a9e8d90845d3e06f1dc4

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
