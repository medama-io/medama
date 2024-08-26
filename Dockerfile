# syntax = docker/dockerfile:1

FROM jdxcode/mise:latest AS build

# Install unzip dependency for bun
# Also combine RUN commands to reduce layers
RUN apt-get update && apt-get install -y unzip && \
	mise install && \
	mise use -g go && \
	mise use -g bun && \
	sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/bin

# Cache build dependencies
ENV GOCACHE=/root/.cache/go-build

# Cache Go modules
WORKDIR /app
COPY core/go.mod core/go.sum ./core/
COPY package.json bun.lockb ./
COPY dashboard/package.json ./dashboard/
COPY tracker/package.json ./tracker/

RUN --mount=type=cache,target=/root/.cache/go-build \
	--mount=type=cache,target=/go/pkg/mod \
	cd core && go mod download && cd .. && \
	bun install --frozen-lockfile

# Copy the rest of the source code
COPY . .
RUN --mount=type=cache,target="/root/.cache/go-build" ~/bin/task core:release

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

