#!/bin/bash

echo "========================================"
echo "Starting all Docker containers..."
echo "========================================"

# Ensure the network exists
echo ""
echo "Setting up Docker network..."
docker network create blog-net 2>/dev/null || echo "Network already exists"

# Start database container
echo ""
echo "Step 1/2: Starting database container..."
./scripts/run_db.sh

# Wait for database to be ready
echo ""
echo "Waiting for database to be ready..."
sleep 3

# Check if database is ready
max_attempts=30
attempt=0
until docker exec blog-postgres pg_isready -U blog -d blogdb > /dev/null 2>&1; do
  attempt=$((attempt + 1))
  if [ $attempt -gt $max_attempts ]; then
    echo "Error: Database failed to start in time"
    exit 1
  fi
  echo "Waiting for database... (attempt $attempt/$max_attempts)"
  sleep 1
done
echo "Database is ready!"

# Start server container
echo ""
echo "Step 2/2: Starting server container..."
./scripts/run_server.sh

echo ""
echo "========================================"
echo "All containers started successfully!"
echo "========================================"
echo ""
echo "Services:"
echo "  - Database: localhost:5432 (blogdb)"
echo "  - Server: http://localhost:8080"
echo ""
echo "To view logs:"
echo "  docker logs -f blog-postgres"
echo "  docker logs -f blog-server"
