# Blog REST API

A simple, production-ready REST API for managing blog posts and comments built with Go and MongoDB.

## Features

- **CRUD Operations** for blog posts and comments
- **MongoDB** integration with proper connection handling
- **Docker** containerization with Docker Compose
- **Clean JSON API** responses with consistent error handling
- **Input validation** and proper HTTP status codes
- **Production-ready** with timeouts and error handling

## Data Models

### BlogPost

- `id` - Unique identifier
- `title` - Post title
- `content` - Post content
- `created` - Creation timestamp
- `comments` - Associated comments (when retrieving single post)

### Comment

- `id` - Unique identifier
- `post_id` - Reference to blog post
- `author` - Comment author name
- `content` - Comment content
- `created` - Creation timestamp

## API Endpoints

| Method | Endpoint                   | Description                            |
| ------ | -------------------------- | -------------------------------------- |
| GET    | `/api/posts`               | Get all blog posts with comment counts |
| POST   | `/api/posts`               | Create a new blog post                 |
| GET    | `/api/posts/{id}`          | Get specific blog post with comments   |
| POST   | `/api/posts/{id}/comments` | Add comment to a blog post             |

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Or Go 1.21+ and MongoDB (for local development)

### Run with Docker

```bash
# Clone the repository
git clone <repository-url>
cd blog-api

# Start the application
docker-compose up -d

# API will be available at http://localhost:8080
```

### Local Development

```bash
# Install dependencies
go mod download

# Set MongoDB connection (optional)
export MONGODB_URI="mongodb://localhost:27017/blog"

# Run the application
go run main.go
```

## API Usage Examples

### Create a Blog Post

```bash
curl -X POST http://localhost:8080/api/posts \
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
curl http://localhost:8080/api/posts
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
curl http://localhost:8080/api/posts/65f1a2b3c4d5e6f7g8h9i0j1
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
curl -X POST http://localhost:8080/api/posts/65f1a2b3c4d5e6f7g8h9i0j1/comments \
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

## Environment Variables

| Variable      | Default                          | Description               |
| ------------- | -------------------------------- | ------------------------- |
| `MONGODB_URI` | `mongodb://localhost:27017/blog` | MongoDB connection string |
| `PORT`        | `8080`                           | Server port               |

## Error Responses

All error responses follow this format:

```json
{
  "success": false,
  "error": "Error message description"
}
```

Common HTTP status codes:

- `400` - Bad Request (invalid JSON, missing required fields)
- `404` - Not Found (post not found)
- `500` - Internal Server Error (database errors)

## Project Structure

```
.
├── app
├──────cmd
├───────── main.go              # Main application file
├── go.mod              # Go module dependencies
├── go.sum              # Go module checksums
├── Dockerfile          # Docker configuration
├── docker-compose.yml  # Docker Compose setup
└── README.md           # This file
```

## Production Considerations

- **Database Indexing**: Add indexes on frequently queried fields
- **Authentication**: Implement JWT or similar auth mechanism
- **Rate Limiting**: Add rate limiting middleware
- **Logging**: Implement structured logging
- **Monitoring**: Add health checks and metrics
- **Validation**: Enhanced input validation and sanitization
- **CORS**: Configure CORS for frontend integration

## Dependencies

- [Fiber/v2](https://docs.gofiber.io/api/fiber) - HTTP router
- [MongoDB Go Driver](https://go.mongodb.org/mongo-driver) - MongoDB client

## License

MIT License
