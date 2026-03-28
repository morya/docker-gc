# Docker GC (Golang Implementation)

A Docker image that allows scheduled cleanup of unused Docker images, containers, and volumes, implemented in Golang.

## Features

- Scheduled garbage collection using cron
- Clean up exited containers older than grace period
- Clean up unused (dangling) images
- Optional cleanup of dangling volumes
- Configurable via environment variables
- Dry run mode for testing

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `CRON` | `0 0 * * *` (midnight daily) | Cron schedule for garbage collection |
| `GRACE_PERIOD_SECONDS` | `3600` (1 hour) | Minimum age for containers/images to be removed |
| `FORCE_CONTAINER_REMOVAL` | `0` | Force removal of containers (set to `1` to enable) |
| `FORCE_IMAGE_REMOVAL` | `0` | Force removal of images (set to `1` to enable) |
| `CLEAN_UP_VOLUMES` | `0` | Clean up dangling volumes (set to `1` to enable) |
| `DRY_RUN` | `0` | Dry run mode - no actual removal (set to `1` to enable) |
| `EXCLUDE_VOLUMES_IDS_FILE` | (none) | Path to file containing volume IDs to exclude |
| `VOLUME_DELETE_ONLY_DRIVER` | (none) | Only delete volumes with specific driver |

## Quick Start

```bash
# Run with default settings (daily at midnight)
docker run -d \
  -v /var/run/docker.sock:/var/run/docker.sock \
  ghcr.io/morya/docker-gc:latest

# Run with custom schedule (every 6 hours)
docker run -d \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e CRON="0 */6 * * *" \
  ghcr.io/morya/docker-gc:latest

# Run with force removal and volume cleanup
docker run -d \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e FORCE_CONTAINER_REMOVAL=1 \
  -e FORCE_IMAGE_REMOVAL=1 \
  -e CLEAN_UP_VOLUMES=1 \
  ghcr.io/morya/docker-gc:latest

# Dry run mode (test without actual removal)
docker run -d \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e DRY_RUN=1 \
  ghcr.io/morya/docker-gc:latest
```

## Building from Source

```bash
# Clone the repository
git clone https://github.com/morya/docker-gc.git
cd docker-gc

# Build the Docker image
docker build -t ghcr.io/morya/docker-gc:latest .

# Run locally
docker run -d \
  -v /var/run/docker.sock:/var/run/docker.sock \
  ghcr.io/morya/docker-gc:latest
```

## Development

```bash
# Build Go binary
go build -o docker-gc-cron ./cmd/docker-gc-cron

# Run locally
./docker-gc-cron
```

## Dockerfile Features

- Two-stage build: golang:alpine for building, alpine for runtime
- CGO_ENABLED=0 for static binary
- Small final image size
- Includes docker-cli for Docker operations
- Multi-architecture support: amd64 and arm64

## GitHub Actions

Automated builds are configured via GitHub Actions. Images are built and pushed to GitHub Container Registry on:
- Push to main/master branch
- Push of tags (v*)
- Pull requests (build only, no push)

## License

Apache License 2.0