# syntax = docker/dockerfile:1

# Adjust NODE_VERSION as desired
FROM node:21-slim

LABEL fly_launch_runtime="Remix"

# Install pnpm
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

COPY . /dashboard
WORKDIR /dashboard

# Install packages needed to build node modules
RUN pnpm install

# Build dashboard
ENV NODE_ENV="production"
ENV API_URL="https://medama-core.fly.dev"
RUN pnpm run build

# Start the server by default, this can be overwritten at runtime
EXPOSE 3000
CMD [ "pnpm", "run", "start" ]
