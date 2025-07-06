// Package models defines the data structures used throughout the blog API.
// It contains request/response models for HTTP endpoints and database entities
// with proper JSON and BSON tags for serialization.
package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreatePostRequest represents the JSON payload for creating a new blog post.
// Used in POST /api/posts endpoint to capture the required fields for post creation.
type CreatePostRequest struct {
	Title   string `json:"title"`   // Post title (required)
	Content string `json:"content"` // Post content/body (required)
}

// CreateCommentRequest represents the JSON payload for creating a new comment.
// Used in POST /api/posts/:id/comments endpoint to capture comment details.
type CreateCommentRequest struct {
	Author  string `json:"author"`  // Comment author name (required)
	Content string `json:"content"` // Comment text content (required)
}

// DeletePostRequest represents the request structure for deleting a blog post.
// Contains the MongoDB ObjectID of the post to be deleted.
// The bson tag supports both JSON requests and direct MongoDB operations.
type DeletePostRequest struct {
	ID primitive.ObjectID `json:"id"` // MongoDB ObjectID of the post
}

// DeleteCommentRequest represents the request structure for deleting a comment.
// Contains the MongoDB ObjectID of the comment to be deleted.
// The bson tag supports both JSON requests and direct MongoDB operations.
type DeleteCommentRequest struct {
	ID primitive.ObjectID `json:"id"` // MongoDB ObjectID of the comment
}

// APIResponse is the standardized response structure for all API endpoints.
// Provides consistent format for success/error responses with optional data payload.
// This ensures uniform client-side response handling across the entire API.
type APIResponse struct {
	Success bool   `json:"success"`         // Indicates if the operation was successful
	Data    any    `json:"data,omitempty"`  // Response payload (omitted if nil/empty)
	Error   string `json:"error,omitempty"` // Error message (omitted if empty)
}
