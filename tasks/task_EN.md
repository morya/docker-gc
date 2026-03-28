# docker-gc Project

Help me reimplement the functionality of `clockworksoul/docker-gc-cron` using Golang.

## Implementation Requirements

- Provide Dockerfile to build the current project
    - Use golang:alpine as the base image for Golang build, implement as a 2-stage build, and use alpine as the base image for the built /app executable.
    - When building the Golang app, set the CGO_ENABLED=0 environment variable in advance for convenient execution.
- Write GitHub Actions configuration files for the current project
- No need to implement test-related unit tests