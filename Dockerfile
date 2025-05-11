# Builder stage - Use Go 1.23.2 to match toolchain
FROM --platform=$BUILDPLATFORM golang:1.23.2-bookworm AS builder

WORKDIR /app

# Set build arguments
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG REACT_APP_API_URL=/api

# Set Go environment
ENV GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    CGO_ENABLED=0 \
    GOAMD64=v2
# Important for compatibility

# Install build dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends git && \
    rm -rf /var/lib/apt/lists/*

# Copy and build
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
# Copy wait-for-it.sh from the API directory itself
COPY ./wait-for-it.sh ./wait-for-it.sh
# Make wait-for-it.sh executable here
RUN chmod +x ./wait-for-it.sh
# Build the application, mounting the cache
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -v -trimpath -ldflags="-s -w" -o /app/main .

# Final stage - Revert to debian-slim to include a shell for wait-for-it.sh
FROM debian:bookworm-slim

# Install bash and netcat needed by the standardized wait-for-it.sh
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash \
    netcat-openbsd \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /root/
# Directories will be created by COPY if they don't exist

COPY --from=builder /app/main .
COPY --from=builder /app/fonts /app/fonts
# Copy the script from the builder stage where it was placed
COPY --from=builder /app/wait-for-it.sh .
# Ensure script is executable (though it should be from builder)
RUN chmod +x ./wait-for-it.sh ./main

# Remove USER directive, run as root (default) for simplicity for now

# Use the shell script as entrypoint
ENTRYPOINT ["./wait-for-it.sh", "db:3306", "--", "./main"]
# Pass arguments via CMD in exec form
CMD ["-dbuser=${MYSQL_USER}", \
     "-dbpassword=${MYSQL_PASSWORD}", \
     "-dbname=${MYSQL_DATABASE}", \
     "-dbhost=${MYSQL_HOST}"]

# NOTE: The MYSQL_USER, MYSQL_PASSWORD, MYSQL_DATABASE, MYSQL_ROOT_PASSWORD,
# and MYSQL_HOST environment variables are expected to be set for Docker Compose.
# These are typically defined in an .env file in the same directory as your
# docker-compose.yml file. For example:
#
# .env file contents:
# MYSQL_USER=youruser
# MYSQL_PASSWORD=yourpassword
# MYSQL_DATABASE=yourdatabase
# MYSQL_ROOT_PASSWORD=yourrootpassword
# MYSQL_HOST=db # Or your MySQL service name in docker-compose.yml
#
# The warnings "variable is not set. Defaulting to a blank string."
# appear when these are not defined in your Docker Compose environment.