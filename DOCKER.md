# Docker Setup Guide

This guide explains how to build and run the Blog API using Docker containers for both the application server and PostgreSQL database.

## Overview

The Blog API Docker setup consists of two main components:

1. **PostgreSQL Database Container** (`blog-postgres`) - Stores blog posts data
2. **Go Application Server Container** (`blog-server`) - Runs the REST API

Both containers communicate through a dedicated Docker network (`blog-net`) and persist data through Docker volumes.

## Prerequisites

- Docker installed and running on your system
- Basic understanding of Docker commands
- Ports 8080 (server) and 5432 (database) available on your host machine
- Minimum 2GB of available disk space for images and volumes

## Architecture

```
┌─────────────────────────────────────────────────┐
│                   Docker Host                    │
│                                                  │
│  ┌────────────────────────────────────────────┐ │
│  │         Docker Network: blog-net           │ │
│  │                                            │ │
│  │  ┌──────────────────┐  ┌────────────────┐ │ │
│  │  │  blog-server     │  │  blog-postgres │ │ │
│  │  │                  │  │                │ │ │
│  │  │  Go Application  │──│  PostgreSQL 16 │ │ │
│  │  │  Port: 8080      │  │  Port: 5432    │ │ │
│  │  └──────────────────┘  └────────────────┘ │ │
│  │                              │             │ │
│  │                              │             │ │
│  │                         ┌────▼────────┐   │ │
│  │                         │ blog_pgdata │   │ │
│  │                         │   (volume)  │   │ │
│  │                         └─────────────┘   │ │
│  └────────────────────────────────────────────┘ │
│                      │                          │
│                      │ Port Mapping             │
└──────────────────────┼──────────────────────────┘
                       │
              ┌────────▼────────┐
              │   Host: 8080    │
              │ (API Accessible)│
              └─────────────────┘
```

## Quick Start

### Step 1: Create Environment File

Create a `.env` file in the project root with the following configuration:

```env
BUILD=prod
PORT=8080
DB_HOST=blog-postgres
DB_PORT=5432
DB_USER=blog
DB_PASSWORD=blogpass
DB_NAME=blogdb
```

**Note:** In production, use strong passwords and consider using Docker secrets.

### Step 2: Build All Components

Build both the database and server images:

```bash
./scripts/build_all.sh
```

This script will:
1. Pull the PostgreSQL 16 image
2. Compile the Go application for Linux ARM64
3. Build the server Docker image with the binary

### Step 3: Start All Services

Start both the database and server containers:

```bash
./scripts/run_all.sh
```

This script will:
1. Create the `blog-net` Docker network (if it doesn't exist)
2. Start the PostgreSQL container with persistent storage
3. Wait for the database to be ready
4. Start the server container connected to the database

### Step 4: Verify Installation

Check that both containers are running:

```bash
docker ps
```

You should see both `blog-server` and `blog-postgres` containers with status "Up".

Test the API:

```bash
curl http://localhost:8080/v1/
```

Expected response:
```json
{"status":"ok"}
```

### Step 5: Stop All Services

When you're done, stop all containers:

```bash
./scripts/stop_all.sh
```

## Individual Component Management

### Database

#### Build (Pull) the Database Image

```bash
./scripts/build_db.sh
```

This pulls the official PostgreSQL 16 image from Docker Hub.

#### Run the Database Container

```bash
./scripts/run_db.sh
```

Database configuration:
- **Container name:** `blog-postgres`
- **Network:** `blog-net`
- **Port mapping:** `5432:5432` (host:container)
- **Database name:** `blogdb`
- **User:** `blog`
- **Password:** `blogpass`
- **Persistent volume:** `blog_pgdata`

#### Connect to Database

Using psql from host:
```bash
psql -h localhost -p 5432 -U blog -d blogdb
```

Using Docker exec:
```bash
docker exec -it blog-postgres psql -U blog -d blogdb
```

### Server

#### Build the Server Image

```bash
./scripts/build_server.sh
```

This script:
1. Compiles the Go binary for Linux ARM64 architecture
2. Builds a Docker image containing:
   - The compiled binary
   - Database migration files
   - Environment configuration

#### Run the Server Container

```bash
./scripts/run_server.sh
```

Server configuration:
- **Container name:** `blog-server`
- **Network:** `blog-net`
- **Port mapping:** `8080:8080` (host:container)
- **Environment:** Loaded from `.env` file
- **Dependencies:** Requires `blog-postgres` to be running

## Script Reference

| Script | Purpose | Prerequisites |
|--------|---------|---------------|
| `build_db.sh` | Pull PostgreSQL image | Docker installed |
| `build_server.sh` | Build server image | Go installed, source files |
| `build_all.sh` | Build all images | Both above |
| `run_db.sh` | Start database | Database image built |
| `run_server.sh` | Start server | Server image built, DB running |
| `run_all.sh` | Start all containers | All images built |
| `stop_all.sh` | Stop and remove containers | None |

## Common Workflows

### Development Workflow

**Initial Setup:**
```bash
# One-time setup
./scripts/build_all.sh
./scripts/run_all.sh
```

**After Code Changes:**
```bash
# Stop the server only
docker stop blog-server
docker rm blog-server

# Rebuild and restart server
./scripts/build_server.sh
./scripts/run_server.sh
```

**Note:** The database keeps running and retains data between server restarts.

### Testing Workflow

```bash
# Start fresh environment
./scripts/stop_all.sh
docker volume rm blog_pgdata
./scripts/build_all.sh
./scripts/run_all.sh

# Run your tests
curl -X POST http://localhost:8080/v1/posts \
  -H "Content-Type: application/json" \
  -d '{"title":"Test","content":"Testing..."}'

# Clean up
./scripts/stop_all.sh
```

### Production Deployment

```bash
# Ensure production environment variables
cat > .env <<EOF
BUILD=prod
PORT=8080
DB_HOST=blog-postgres
DB_PORT=5432
DB_USER=blog
DB_PASSWORD=your-secure-password
DB_NAME=blogdb
EOF

# Build and deploy
./scripts/build_all.sh
./scripts/run_all.sh

# Monitor logs
docker logs -f blog-server
```

### Clean Start

```bash
# Complete cleanup and restart
./scripts/stop_all.sh
docker volume rm blog_pgdata
docker network rm blog-net
./scripts/build_all.sh
./scripts/run_all.sh
```

## Monitoring and Logs

### View Real-Time Logs

**Server logs:**
```bash
docker logs -f blog-server
```

**Database logs:**
```bash
docker logs -f blog-postgres
```

**Both logs simultaneously:**
```bash
docker logs -f blog-server & docker logs -f blog-postgres
```

### Container Status

**Check running containers:**
```bash
docker ps
```

**Check all containers (including stopped):**
```bash
docker ps -a
```

**Inspect container details:**
```bash
docker inspect blog-server
docker inspect blog-postgres
```

### Resource Usage

**View resource consumption:**
```bash
docker stats blog-server blog-postgres
```

### Network Information

**Inspect the blog network:**
```bash
docker network inspect blog-net
```

## Data Management

### Database Backups

**Create a backup:**
```bash
docker exec blog-postgres pg_dump -U blog blogdb > backup.sql
```

**Restore from backup:**
```bash
cat backup.sql | docker exec -i blog-postgres psql -U blog -d blogdb
```

**Export backup from container:**
```bash
docker exec blog-postgres pg_dump -U blog blogdb -f /tmp/backup.sql
docker cp blog-postgres:/tmp/backup.sql ./backup.sql
```

### Volume Management

**List volumes:**
```bash
docker volume ls
```

**Inspect volume:**
```bash
docker volume inspect blog_pgdata
```

**Backup volume data:**
```bash
docker run --rm -v blog_pgdata:/data -v $(pwd):/backup ubuntu tar czf /backup/blog_pgdata_backup.tar.gz /data
```

**Restore volume data:**
```bash
docker run --rm -v blog_pgdata:/data -v $(pwd):/backup ubuntu tar xzf /backup/blog_pgdata_backup.tar.gz -C /
```

## Cleanup

### Remove Containers Only

```bash
./scripts/stop_all.sh
```

### Remove Network

```bash
docker network rm blog-net
```

### Remove Persistent Database Data

**WARNING: This will delete all blog posts!**

```bash
docker volume rm blog_pgdata
```

### Complete Cleanup

Remove all components including images:

```bash
# Stop and remove containers
./scripts/stop_all.sh

# Remove network
docker network rm blog-net

# Remove volume (deletes data!)
docker volume rm blog_pgdata

# Remove images
docker rmi blog-server:latest
docker rmi postgres:16
```

### Prune Unused Resources

```bash
# Remove all stopped containers
docker container prune -f

# Remove all unused volumes
docker volume prune -f

# Remove all unused networks
docker network prune -f

# Remove all unused images
docker image prune -a -f
```

## Troubleshooting

### Port Already in Use

**Problem:** Port 8080 or 5432 is already in use.

**Solution:**
```bash
# Find process using the port
lsof -i :8080
lsof -i :5432

# Kill the process or stop the service
# Then restart the containers
./scripts/run_all.sh
```

**Alternative:** Modify port mappings in the run scripts.

### Container Won't Start

**Problem:** Container exits immediately after starting.

**Check logs:**
```bash
docker logs blog-server
docker logs blog-postgres
```

**Common causes:**
- Environment variables not set correctly
- Database not ready when server starts
- Port conflicts
- Insufficient resources

**Solution:**
1. Verify `.env` file exists and has correct values
2. Check Docker logs for specific errors
3. Ensure database starts before server
4. Verify port availability

### Database Connection Issues

**Problem:** Server can't connect to database.

**Verify database is running:**
```bash
docker ps | grep blog-postgres
```

**Check network connectivity:**
```bash
docker network inspect blog-net
```

**Verify credentials:**
1. Check `.env` file matches database configuration
2. Try connecting manually:
```bash
docker exec -it blog-postgres psql -U blog -d blogdb
```

**Solution:**
```bash
# Restart both containers
./scripts/stop_all.sh
./scripts/run_all.sh
```

### Network Issues

**Problem:** Containers can't communicate.

**Recreate network:**
```bash
docker network rm blog-net
docker network create blog-net
./scripts/run_all.sh
```

**Verify network:**
```bash
docker network inspect blog-net
```

Both containers should appear in the "Containers" section.

### Build Failures

**Problem:** Image build fails.

**For database:**
```bash
# Check internet connectivity
docker pull postgres:16
```

**For server:**
```bash
# Verify Go compilation works
go build -o blog-server

# Check Dockerfile syntax
docker build -t blog-server:latest .
```

### Permission Denied Errors

**Problem:** Permission denied when running scripts.

**Solution:**
```bash
chmod +x scripts/*.sh
```

### Out of Disk Space

**Problem:** No space left on device.

**Check Docker disk usage:**
```bash
docker system df
```

**Clean up:**
```bash
docker system prune -a --volumes
```

## Advanced Configuration

### Custom Database Configuration

Edit `scripts/run_db.sh` to add PostgreSQL parameters:

```bash
docker run -d \
  --name blog-postgres \
  --network blog-net \
  -p 5432:5432 \
  -e POSTGRES_USER=blog \
  -e POSTGRES_PASSWORD=blogpass \
  -e POSTGRES_DB=blogdb \
  -c max_connections=200 \
  -c shared_buffers=256MB \
  -v blog_pgdata:/var/lib/postgresql/data \
  postgres:16
```

### Environment-Specific Builds

**Development with hot reload:**
```bash
docker run -d \
  --name blog-server-dev \
  --network blog-net \
  -p 8080:8080 \
  -v $(pwd):/app \
  -w /app \
  --env-file .env \
  golang:1.21 \
  go run main.go
```

### Docker Compose Alternative

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  db:
    image: postgres:16
    container_name: blog-postgres
    environment:
      POSTGRES_USER: blog
      POSTGRES_PASSWORD: blogpass
      POSTGRES_DB: blogdb
    volumes:
      - blog_pgdata:/var/lib/postgresql/data
    networks:
      - blog-net
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U blog"]
      interval: 5s
      timeout: 5s
      retries: 5

  server:
    build: .
    container_name: blog-server
    env_file: .env
    depends_on:
      db:
        condition: service_healthy
    networks:
      - blog-net
    ports:
      - "8080:8080"

networks:
  blog-net:
    driver: bridge

volumes:
  blog_pgdata:
```

**Usage:**
```bash
docker-compose up -d
docker-compose down
docker-compose logs -f
```

## Security Best Practices

1. **Use Strong Passwords:** Change default passwords in `.env`
2. **Don't Commit Secrets:** Add `.env` to `.gitignore`
3. **Use Docker Secrets:** For production deployments
4. **Limit Container Permissions:** Run containers as non-root user
5. **Keep Images Updated:** Regularly update base images
6. **Network Isolation:** Use dedicated networks for different services
7. **Volume Encryption:** Consider encrypting sensitive data volumes

## Performance Optimization

1. **Multi-stage Builds:** Reduce image size
2. **Layer Caching:** Organize Dockerfile to maximize cache hits
3. **Resource Limits:** Set memory and CPU limits
4. **Connection Pooling:** Configure database connection pool
5. **Health Checks:** Implement proper health check endpoints

## Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [PostgreSQL Docker Image](https://hub.docker.com/_/postgres)
- [Go Docker Best Practices](https://docs.docker.com/language/golang/)
- [Project README](README.md) for API documentation

## Support

For issues related to Docker setup:
1. Check this troubleshooting guide
2. Review container logs
3. Verify environment configuration
4. Open an issue on GitHub with error details
