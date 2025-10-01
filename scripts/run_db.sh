#!/bin/bash

# Stop and remove existing database container if it exists
echo "Stopping and removing existing database container..."
docker stop blog-postgres 2>/dev/null || true
docker rm blog-postgres 2>/dev/null || true

# Ensure the network exists
echo "Ensuring Docker network exists..."
docker network create blog-net 2>/dev/null || true

# Run the database container
echo "Starting database container..."
docker run -d \
  --name blog-postgres \
  --network blog-net \
  -e POSTGRES_USER=blog \
  -e POSTGRES_PASSWORD=blogpass \
  -e POSTGRES_DB=blogdb \
  -v blog_pgdata:/var/lib/postgresql/data \
  -p 5432:5432 \
  postgres:16

echo "Database container started successfully!"
echo "PostgreSQL is running on localhost:5432"
echo "Database: blogdb, User: blog"
