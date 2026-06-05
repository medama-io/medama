# syntax = docker/dockerfile:1@sha256:87999aa3d42bdc6bea60565083ee17e86d1f3339802f543c0d03998580f9cb89

ARG MANYLINUX_IMAGE=quay.io/pypa/manylinux_2_34_x86_64@sha256:42b9c13a7cb538e59c2faf75835e8f7dd3c7a764ed7546937b71e2cfec980d6e
FROM ${MANYLINUX_IMAGE} AS build

ARG VERSION=development
ARG COMMIT_SHA=development

ENV VERSION=${VERSION}
ENV COMMIT_SHA=${COMMIT_SHA}

RUN yum -y install git curl unzip zip

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
ENV MISE_DATA_DIR="/mise"
ENV MISE_CONFIG_DIR="/mise"
ENV MISE_CACHE_DIR="/mise/cache"
ENV MISE_INSTALL_PATH="/usr/local/bin/mise"
ENV MISE_EXPERIMENTAL="1"
ENV PATH="/mise/shims:$PATH"

RUN curl https://mise.run | sh

WORKDIR /app

# Install language runtimes from mise.toml
COPY mise.toml ./mise.toml
COPY core/mise.toml ./core/mise.toml
COPY dashboard/mise.toml ./dashboard/mise.toml
COPY tracker/mise.toml ./tracker/mise.toml
RUN mise trust -a -y && mise install go bun

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
RUN --mount=type=cache,target=${GOCACHE} mise run //core:release:docker

# Build the final image
FROM gcr.io/distroless/cc-debian12@sha256:aa0b7af67fa8211751ea6e00baa8373ba56cc1417ffc986ec9619bd0e1556b56

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
