#!/bin/bash

# Stop and remove existing server container if it exists
echo "Stopping and removing existing server container..."
docker stop blog-server 2>/dev/null || true
docker rm blog-server 2>/dev/null || true

# Ensure the network exists
echo "Ensuring Docker network exists..."
docker network create blog-net 2>/dev/null || true

# Run the server container
echo "Starting server container..."
docker run -d \
  --name blog-server \
  --network blog-net \
  --env-file .env \
  -p 8080:8080 \
  blog-server:latest

echo "Server container started successfully!"
echo "Server is running at http://localhost:8080"
