#!/bin/bash

# Set environment variables for cross-compilation
export GOOS=linux
export GOARCH=amd64

# Build the Go application
go build -o timengledev_blog ./cmd/web

# Build the Docker image
docker build --platform=linux/amd64 -t timenglesf/timengledev_blog:latest .
