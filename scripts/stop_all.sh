#!/bin/bash

echo "========================================"
echo "Stopping all Docker containers..."
echo "========================================"

# Stop server container
echo ""
echo "Stopping server container..."
docker stop blog-server 2>/dev/null || echo "Server container not running"
docker rm blog-server 2>/dev/null || echo "Server container already removed"

# Stop database container
echo ""
echo "Stopping database container..."
docker stop blog-postgres 2>/dev/null || echo "Database container not running"
docker rm blog-postgres 2>/dev/null || echo "Database container already removed"

echo ""
echo "========================================"
echo "All containers stopped successfully!"
echo "========================================"
echo ""
echo "Note: The Docker network 'blog-net' and volume 'blog_pgdata' are preserved."
echo "To remove them, run:"
echo "  docker network rm blog-net"
echo "  docker volume rm blog_pgdata"
