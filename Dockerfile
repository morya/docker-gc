# Build stage
FROM golang:alpine AS builder

WORKDIR /app

# Set environment variables for build
ARG TARGETARCH
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=${TARGETARCH}

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY cmd/ ./cmd/

# Build the application
RUN go build -o docker-gc-cron ./cmd/docker-gc-cron

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install necessary packages
RUN apk --no-cache add docker-cli bash tzdata

# Copy the binary from builder
COPY --from=builder /app/docker-gc-cron .

# Create necessary directories
RUN mkdir -p /var/log

# Set up entrypoint
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Default environment variables
ENV CRON="0 0 * * *"
ENV GRACE_PERIOD_SECONDS=3600
ENV MINIMUM_IMAGES_TO_SAVE=0
ENV FORCE_CONTAINER_REMOVAL=0
ENV FORCE_IMAGE_REMOVAL=0
ENV CLEAN_UP_VOLUMES=0
ENV DRY_RUN=0

# Mount Docker socket
VOLUME /var/run/docker.sock

ENTRYPOINT ["/entrypoint.sh"]