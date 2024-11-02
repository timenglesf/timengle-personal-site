#!/bin/bash

# Set environment variables for cross-compilation
export GOOS=linux
export GOARCH=arm64

# Build the Go application
go build -o timengledev_blog ./cmd/web
