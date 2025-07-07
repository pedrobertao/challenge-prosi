# Blog API Request Documentation

This document outlines all the HTTP requests supported by the blog API, including endpoints, request formats, and expected responses.

## Base URL

```
https://challenge-prosi-390352505094.southamerica-east1.run.app/api
```

## Authentication

No authentication required for any endpoints.

---

## Posts Endpoints

### 1. Get All Posts

**Endpoint:** `GET /api/posts`

**Description:** Retrieves a list of all blog posts with summary information including post ID, title, comment count, and creation date.

**Request:**

```http
GET /api/posts
Content-Type: application/json
```

**Response Examples:**

**Success (200):**

```json
{
  "success": true,
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "title": "My First Blog Post",
      "comment_count": 5,
      "created_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "507f1f77bcf86cd799439012",
      "title": "Another Great Post",
      "comment_count": 2,
      "created_at": "2024-01-16T14:22:00Z"
    }
  ],
  "error": ""
}
```

**No Posts Found (404):**

```json
{
  "success": true,
  "data": [],
  "error": ""
}
```

**Database Error (502):**

```json
{
  "success": false,
  "error": "Failed to fetch posts"
}
```

---

### 2. Create New Post

**Endpoint:** `POST /api/posts`

**Description:** Creates a new blog post with the provided title and content.

**Request:**

```http
POST /api/posts
Content-Type: application/json

{
  "title": "My New Blog Post",
  "content": "This is the content of my new blog post. It can be quite long and contain multiple paragraphs."
}
```

**Response Examples:**

**Success (200):**

```json
{
  "success": true,
  "data": {
    "id": "507f1f77bcf86cd799439013",
    "title": "My New Blog Post",
    "content": "This is the content of my new blog post. It can be quite long and contain multiple paragraphs.",
    "created_at": "2024-01-17T09:15:00Z",
    "comments": []
  },
  "error": ""
}
```

**Invalid JSON (400):**

```json
{
  "success": false,
  "error": "Invalid JSON"
}
```

**Missing Required Fields (400):**

```json
{
  "success": false,
  "error": "Title and content required"
}
```

**Database Error (502):**

```json
{
  "success": false,
  "error": "Failed to create post"
}
```

---

### 3. Get Single Post

**Endpoint:** `GET /api/posts/:id`

**Description:** Retrieves a specific blog post by its ID along with all associated comments.

**Request:**

```http
GET /api/posts/507f1f77bcf86cd799439011
Content-Type: application/json
```

**Response Examples:**

**Success (200):**

```json
{
  "success": true,
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "title": "My First Blog Post",
    "content": "This is the full content of my first blog post with all the details.",
    "created_at": "2024-01-15T10:30:00Z",
    "comments": [
      {
        "id": "507f1f77bcf86cd799439021",
        "post_id": "507f1f77bcf86cd799439011",
        "author": "John Doe",
        "content": "Great post! Thanks for sharing.",
        "created_at": "2024-01-15T12:45:00Z"
      },
      {
        "id": "507f1f77bcf86cd799439022",
        "post_id": "507f1f77bcf86cd799439011",
        "author": "Jane Smith",
        "content": "I found this very helpful.",
        "created_at": "2024-01-15T14:20:00Z"
      }
    ]
  },
  "error": ""
}
```

**Invalid Post ID (400):**

```json
{
  "success": false,
  "error": "Invalid post ID"
}
```

**Post Not Found (404):**

```json
{
  "success": false,
  "error": "Post not found"
}
```

**Database Error (500):**

```json
{
  "success": false,
  "error": "Failed to fetch post"
}
```

---

### 4. Delete Post

**Endpoint:** `DELETE /api/posts/:id`

**Description:** Deletes a specific blog post and all its associated comments atomically using MongoDB transactions.

**Request:**

```http
DELETE /api/posts/507f1f77bcf86cd799439011
Content-Type: application/json
```

**Response Examples:**

**Success (200):**

```json
{
  "success": true,
  "data": "507f1f77bcf86cd799439011",
  "error": ""
}
```

**Invalid Post ID (400):**

```json
{
  "success": false,
  "error": "Invalid post ID"
}
```

**Post Not Found (400):**

```json
{
  "success": false,
  "error": "Post not found"
}
```

**Database Error (502):**

```json
{
  "success": false,
  "error": "Failed to delete post"
}
```

---

## Comments Endpoints

### 5. Create Comment

**Endpoint:** `POST /api/posts/:id/comments`

**Description:** Creates a new comment on a specific blog post.

**Request:**

```http
POST /api/posts/507f1f77bcf86cd799439011/comments
Content-Type: application/json

{
  "author": "John Doe",
  "content": "This is a great blog post! Thanks for sharing your insights."
}
```

**Response Examples:**

**Success (200):**

```json
{
  "success": true,
  "data": {
    "id": "507f1f77bcf86cd799439023",
    "post_id": "507f1f77bcf86cd799439011",
    "author": "John Doe",
    "content": "This is a great blog post! Thanks for sharing your insights.",
    "created_at": "2024-01-17T11:30:00Z"
  },
  "error": ""
}
```

**Invalid Post ID (400):**

```json
{
  "success": false,
  "error": "Invalid post ID"
}
```

**Invalid JSON (400):**

```json
{
  "success": false,
  "error": "Invalid JSON"
}
```

**Missing Required Fields (400):**

```json
{
  "success": false,
  "error": "Author and content required"
}
```

**Post Not Found (404):**

```json
{
  "success": false,
  "error": "Post not found"
}
```

**Database Error (500):**

```json
{
  "success": false,
  "error": "Failed to create comment"
}
```

---

### 6. Delete Comment

**Endpoint:** `DELETE /api/comments/:id`

**Description:** Deletes a specific comment by its ID.

**Request:**

```http
DELETE /api/comments/507f1f77bcf86cd799439023
Content-Type: application/json
```

**Response Examples:**

**Success (200):**

```json
{
  "success": true,
  "data": "507f1f77bcf86cd799439023",
  "error": ""
}
```

**Invalid Comment ID (400):**

```json
{
  "success": false,
  "error": "Invalid comment ID"
}
```

**Comment Not Found (400):**

```json
{
  "success": false,
  "error": "No comment found to delete"
}
```

**Database Error (502):**

```json
{
  "success": false,
  "error": "Failed to delete comment"
}
```

---

## Request/Response Format

### Common Response Structure

All API responses follow this consistent format:

```json
{
  "success": boolean,
  "data": object | array | string | null,
  "error": string
}
```

### ID Format

All IDs in the API are MongoDB ObjectIDs represented as 24-character hexadecimal strings.

**Example:** `507f1f77bcf86cd799439011`

### Date Format

All timestamps are in ISO 8601 format with UTC timezone.

**Example:** `2024-01-17T11:30:00Z`

### Content-Type

All requests should include the `Content-Type: application/json` header when sending JSON data.

---

## Error Handling

The API uses standard HTTP status codes:

- **200**: Success
- **400**: Bad Request (invalid data, missing fields, invalid ID format)
- **404**: Not Found (post or comment doesn't exist)
- **500**: Internal Server Error (database query errors)
- **502**: Bad Gateway (database connection or transaction errors)

All error responses include a descriptive error message in the `error` field and set `success` to `false`.

---

## Database Operations

- **Timeouts**: All database operations have a 10-second timeout to prevent hanging requests
- **Transactions**: Post deletion uses MongoDB transactions to ensure atomicity
- **Validation**: All ObjectIDs are validated before database operations
- **Error Logging**: Database errors are logged with structured logging using Zap
