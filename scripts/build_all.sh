#!/bin/bash

echo "========================================"
echo "Building all Docker images..."
echo "========================================"

# Build database image (pull PostgreSQL)
echo ""
echo "Step 1/2: Building database image..."
./scripts/build_db.sh

# Build server image
echo ""
echo "Step 2/2: Building server image..."
./scripts/build_server.sh

echo ""
echo "========================================"
echo "All images built successfully!"
echo "========================================"
