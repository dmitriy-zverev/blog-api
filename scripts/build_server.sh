#!/bin/bash

# Build the Go binary for Linux ARM64
echo "Building Go binary..."
GOOS=linux GOARCH=arm64 go build -o blog-server .

# Build the Docker image
echo "Building server Docker image..."
docker build --platform linux/arm64 -t blog-server:latest .

echo "Server image built successfully!"
