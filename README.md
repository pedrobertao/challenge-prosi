# Blog REST API

A simple, production-ready REST API for managing blog posts and comments built with **Go Fiber v2** and **MongoDB**.

## Features

- **CRUD Operations** for blog posts and comments
- **MongoDB** integration with proper connection handling
- **Docker** containerization with Docker Compose
- **Clean JSON API** responses with consistent error handling
- **Input validation** and proper HTTP status codes
- **Production-ready** with timeouts and graceful shutdown
- **Fiber v2** - Fast, Express-inspired web framework

## Tech Stack

- **Go 1.24** with Fiber v2 framework
- **MongoDB 7.0** for data storage
- **Docker & Docker Compose** for containerization
- **Alpine Linux** for minimal container size

## Project Structure

```
blog-api/
├── app/
│   └── cmd/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── models/
│   │   └── models.go        # Data models
│   ├── handlers/
│   │   └── handlers.go      # HTTP handlers
│   ├── storage/
│   │   └── storage.go       # Database connection
│   └── routes/
│       └── routes.go        # Route setup
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── .env
├── .env.example
└── README.md
```

## Data Models

### BlogPost

- `id` - Unique identifier (MongoDB ObjectID)
- `title` - Post title
- `content` - Post content
- `created` - Creation timestamp
- `comments` - Associated comments (when retrieving single post)

### Comment

- `id` - Unique identifier (MongoDB ObjectID)
- `post_id` - Reference to blog post
- `author` - Comment author name
- `content` - Comment content
- `created` - Creation timestamp

## API Endpoints

| Method | Endpoint                  | Description                            |
| ------ | ------------------------- | -------------------------------------- |
| GET    | `/api/posts`              | Get all blog posts with comment counts |
| POST   | `/api/posts`              | Create a new blog post                 |
| GET    | `/api/posts/:id`          | Get specific blog post with comments   |
| POST   | `/api/posts/:id/comments` | Add comment to a blog post             |

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Or Go 1.24+ and MongoDB (for local development)

### Run with Docker (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd blog-api

# Create environment file
cp .env.example .env

# Start the application
docker-compose up -d

# API will be available at http://localhost:3030
```

### Local Development

```bash
# Install dependencies
go mod download

# Set environment variables
export MONGODB_URI="mongodb://localhost:27017"
export DB_NAME="blog"
export PORT="3030"

# Run the application
go run app/cmd/main.go
```

## Environment Variables

Create a `.env` file in the project root:

```bash
# Database Configuration
MONGODB_URI=mongodb://mongodb:27017
DB_NAME=blog

# Server Configuration
PORT=3030

# Environment
ENV=production
```

## API Usage Examples

### Create a Blog Post

```bash
curl -X POST http://localhost:3030/api/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Blog Post",
    "content": "This is the content of my first blog post."
  }'
```

**Response:**

```json
{
  "success": true,
  "data": {
    "id": "65f1a2b3c4d5e6f7g8h9i0j1",
    "title": "My First Blog Post",
    "content": "This is the content of my first blog post.",
    "created": "2024-03-15T10:30:00Z"
  }
}
```

### Get All Blog Posts

```bash
curl http://localhost:3030/api/posts
```

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": "65f1a2b3c4d5e6f7g8h9i0j1",
      "title": "My First Blog Post",
      "comment_count": 2,
      "created": "2024-03-15T10:30:00Z"
    }
  ]
}
```

### Get Specific Blog Post

```bash
curl http://localhost:3030/api/posts/65f1a2b3c4d5e6f7g8h9i0j1
```

**Response:**

```json
{
  "success": true,
  "data": {
    "id": "65f1a2b3c4d5e6f7g8h9i0j1",
    "title": "My First Blog Post",
    "content": "This is the content of my first blog post.",
    "comments": [
      {
        "id": "65f1a2b3c4d5e6f7g8h9i0j2",
        "post_id": "65f1a2b3c4d5e6f7g8h9i0j1",
        "author": "John Doe",
        "content": "Great post!",
        "created": "2024-03-15T11:00:00Z"
      }
    ],
    "created": "2024-03-15T10:30:00Z"
  }
}
```

### Add Comment to Blog Post

```bash
curl -X POST http://localhost:3030/api/posts/65f1a2b3c4d5e6f7g8h9i0j1/comments \
  -H "Content-Type: application/json" \
  -d '{
    "author": "Jane Smith",
    "content": "Thanks for sharing this!"
  }'
```

**Response:**

```json
{
  "success": true,
  "data": {
    "id": "65f1a2b3c4d5e6f7g8h9i0j3",
    "post_id": "65f1a2b3c4d5e6f7g8h9i0j1",
    "author": "Jane Smith",
    "content": "Thanks for sharing this!",
    "created": "2024-03-15T11:15:00Z"
  }
}
```

## API Response Format

All API responses follow this consistent structure:

### Success Response:

```json
{
  "success": true,
  "data": { ... }
}
```

### Error Response:

```json
{
  "success": false,
  "error": "Error message description"
}
```

## Error Responses

Common HTTP status codes:

- `400` - Bad Request (invalid JSON, missing required fields)
- `404` - Not Found (post not found)
- `500` - Internal Server Error (database errors)

## Docker Commands

```bash
# Start services
docker-compose up -d

# View logs
docker-compose logs -f blog-api

# Stop services
docker-compose down

# Rebuild after code changes
docker-compose up --build

# Remove volumes (reset database)
docker-compose down -v
```

## Development

### Project Structure Benefits

- **Clean Architecture** - Separation of concerns
- **Testable** - Easy to unit test individual components
- **Maintainable** - Clear organization
- **Scalable** - Easy to add new features
- **Standard** - Follows Go project layout conventions

### Adding New Features

1. Add new models in `internal/models/`
2. Create handlers in `internal/handlers/`
3. Register routes in `internal/routes/`
4. Update documentation

## Production Considerations

- **Database Indexing** - Add indexes on frequently queried fields
- **Authentication** - Implement JWT or similar auth mechanism
- **Rate Limiting** - Add rate limiting middleware
- **Logging** - Implement structured logging with logrus or zap
- **Monitoring** - Add health checks and metrics
- **Validation** - Enhanced input validation and sanitization
- **CORS** - Configure CORS for frontend integration
- **SSL/TLS** - Enable HTTPS in production
- **Load Balancing** - Use reverse proxy (nginx) for multiple instances

## Dependencies

- [Fiber v2](https://github.com/gofiber/fiber) - HTTP web framework
- [MongoDB Go Driver](https://go.mongodb.org/mongo-driver) - MongoDB client
- [Docker](https://www.docker.com/) - Containerization
- [Docker Compose](https://docs.docker.com/compose/) - Multi-container orchestration
