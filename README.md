# Blog API

A RESTful API for a personal blogging platform built with Go and PostgreSQL.

## Features

The API allows users to perform the following operations:

- Create a new blog post
- Update an existing blog post
- Delete an existing blog post
- Get a single blog post
- Get all blog posts
- Filter blog posts by a search term

## Tech Stack

- **Language**: Go 1.21+
- **Database**: PostgreSQL 16
- **SQL Generator**: sqlc
- **Environment Management**: godotenv
- **Containerization**: Docker

## Project Structure

```
blog-api/
├── internal/
│   ├── db/              # Generated database queries (sqlc)
│   ├── handlers/        # HTTP request handlers
│   └── models/          # Data models
├── scripts/             # Build and run scripts
├── sql/
│   ├── queries/         # SQL queries for sqlc
│   └── schema/          # Database schema migrations
├── main.go              # Application entry point
├── routes.go            # Route definitions
├── Dockerfile           # Docker configuration
├── .env                 # Environment variables (create this)
└── sqlc.yaml            # sqlc configuration
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 16 or higher (or use Docker)
- sqlc (optional, for regenerating queries)

## Quick Start

### Option 1: Using Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/dmitriy-zverev/blog-api.git
   cd blog-api
   ```

2. **Create environment file**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Build and run with Docker**
   ```bash
   ./scripts/build_all.sh
   ./scripts/run_all.sh
   ```

4. **The API will be available at**
   ```
   http://localhost:8080
   ```

See [DOCKER.md](DOCKER.md) for detailed Docker instructions.

### Option 2: Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/dmitriy-zverev/blog-api.git
   cd blog-api
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up PostgreSQL database**
   ```bash
   # Create database
   createdb blogdb
   
   # Run migrations
   psql -d blogdb -f sql/schema/001_posts.sql
   ```

4. **Create .env file**
   ```env
   BUILD=dev
   PORT=8080
   DB_URL=postgres://your_user:your_password@localhost:5432/blogdb?sslmode=disable
   ```

5. **Run the server**
   ```bash
   go run .
   ```

## Environment Variables

### Development Mode (`BUILD=dev`)
```env
BUILD=dev
PORT=8080
DB_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
```

### Production Mode (`BUILD=prod`)
```env
BUILD=prod
PORT=8080
DB_HOST=blog-postgres
DB_PORT=5432
DB_USER=blog
DB_PASSWORD=blogpass
DB_NAME=blogdb
```

## API Endpoints

All endpoints are prefixed with `/v1`.

### Base Endpoint

#### Check API Status
```http
GET /v1/
```

**Response:**
```json
OK
```

### Posts Endpoints

#### Create a New Post
```http
POST /v1/posts
Content-Type: application/json

{
  "title": "My First Blog Post",
  "content": "This is the content of my first blog post.",
  "category": "Technology",
  "tags": ["programming", "tutorial", "beginner"]
}
```

**Required Fields:**
- `title` (string) - The title of the blog post
- `content` (string) - The main content of the blog post
- `category` (string) - The category of the blog post
- `tags` (array of strings) - Tags associated with the blog post

**Response (201 Created):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "My First Blog Post",
  "content": "This is the content of my first blog post.",
  "category": "Technology",
  "tags": ["programming", "tutorial", "beginner"],
  "created_at": "2025-01-10T15:00:00Z",
  "updated_at": "2025-01-10T15:00:00Z"
}
```

#### Get All Posts
```http
GET /v1/posts
```

**Optional Query Parameters:**
- `term` - Search term to filter posts by title or content

**Example:**
```http
GET /v1/posts?term=blog
```

**Response (200 OK):**
```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "title": "My First Blog Post",
    "content": "This is the content of my first blog post.",
    "category": "Technology",
    "tags": ["programming", "tutorial", "beginner"],
    "created_at": "2025-01-10T15:00:00Z",
    "updated_at": "2025-01-10T15:00:00Z"
  }
]
```

#### Get a Single Post
```http
GET /v1/posts/{postId}
```

**Response (200 OK):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "My First Blog Post",
  "content": "This is the content of my first blog post.",
  "category": "Technology",
  "tags": ["programming", "tutorial", "beginner"],
  "created_at": "2025-01-10T15:00:00Z",
  "updated_at": "2025-01-10T15:00:00Z"
}
```

#### Update a Post
```http
PUT /v1/posts/{postId}
Content-Type: application/json

{
  "title": "Updated Title",
  "content": "Updated content.",
  "category": "Programming",
  "tags": ["golang", "api", "web-development"]
}
```

**Required Fields:**
- `title` (string) - The title of the blog post
- `content` (string) - The main content of the blog post
- `category` (string) - The category of the blog post
- `tags` (array of strings) - Tags associated with the blog post

**Response (200 OK):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "Updated Title",
  "content": "Updated content.",
  "category": "Programming",
  "tags": ["golang", "api", "web-development"],
  "created_at": "2025-01-10T15:00:00Z",
  "updated_at": "2025-01-10T15:30:00Z"
}
```

#### Delete a Post
```http
DELETE /v1/posts/{postId}
```

**Response (204 No Content)**

## Error Responses

The API uses standard HTTP status codes:

- `200 OK` - Request succeeded
- `201 Created` - Resource created successfully
- `204 No Content` - Request succeeded with no content to return
- `400 Bad Request` - Invalid request body or parameters
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

**Error Response Format:**
```json
{
  "error": "Error message describing what went wrong"
}
```

## Development

### Generating Database Queries
If you modify SQL queries in `sql/queries/`, regenerate the Go code:
```bash
sqlc generate
```

### Building the Binary
```bash
go build -o blog-server
```

### Database Migrations
To add a new migration, create a file in `sql/schema/` with a numeric prefix:
```bash
touch sql/schema/002_add_new_table.sql
```

## Scripts

The `scripts/` directory contains helper scripts for Docker operations:

- `build_all.sh` - Build both database and server images
- `run_all.sh` - Start both database and server containers
- `stop_all.sh` - Stop and remove all containers
- `build_db.sh` - Pull PostgreSQL Docker image
- `build_server.sh` - Build server Docker image
- `run_db.sh` - Start database container
- `run_server.sh` - Start server container
- `start_server.sh` - Start server in local development mode

## Examples

### Using curl

**Create a post:**
```bash
curl -X POST http://localhost:8080/v1/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Hello World",
    "content": "This is my first post!",
    "category": "General",
    "tags": ["introduction", "first-post"]
  }'
```

**Get all posts:**
```bash
curl http://localhost:8080/v1/posts
```

**Search posts:**
```bash
curl "http://localhost:8080/v1/posts?term=hello"
```

**Get a specific post:**
```bash
curl http://localhost:8080/v1/posts/{post-id}
```

**Update a post:**
```bash
curl -X PUT http://localhost:8080/v1/posts/{post-id} \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Title",
    "content": "Updated content",
    "category": "Updated Category",
    "tags": ["updated", "modified"]
  }'
```

**Delete a post:**
```bash
curl -X DELETE http://localhost:8080/v1/posts/{post-id}
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Support

For issues, questions, or contributions, please open an issue on GitHub.

Inspired by [Roadmap.sh](https://roadmap.sh/projects/blogging-platform-api)