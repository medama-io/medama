# syntax = docker/dockerfile:1@sha256:b6afd42430b15f2d2a4c5a02b919e98a525b785b1aaff16747d2f623364e39b6

ARG MANYLINUX_IMAGE=quay.io/pypa/manylinux_2_34_x86_64@sha256:f4cf0b6ca26463dbaebdcc92af683ef6f37d710facf6c1914b30bd69fa2b4d5a
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
FROM gcr.io/distroless/cc-debian12@sha256:0c8eac8ea42a167255d03c3ba6dfad2989c15427ed93d16c53ef9706ea4691df

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
