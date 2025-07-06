package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BlogPost represents a blog post entity stored in MongoDB.
// Contains the full post data including metadata and associated comments.
type BlogPost struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`      // MongoDB ObjectID
	Title     string             `json:"title" bson:"title"`           // Post title
	Content   string             `json:"content" bson:"content"`       // Post content/body
	CreatedAt time.Time          `json:"created_at" bson:"created_at"` // Creation timestamp
	Comments  []Comment          `json:"comments,omitempty" bson:"-"`  // Associated comments (not stored in post document)
}

// BlogPostSummary represents a condensed view of a blog post for list endpoints.
// Used in GET /api/posts to provide overview information without full content.
// Optimized for performance by excluding the potentially large content field.
type BlogPostSummary struct {
	ID           primitive.ObjectID `json:"id"`            // MongoDB ObjectID
	Title        string             `json:"title"`         // Post title
	CommentCount int64              `json:"comment_count"` // Number of comments on this post
	CreatedAt    time.Time          `json:"created_at"`    // Creation timestamp
}

// Comment represents a comment entity stored in MongoDB.
// Comments are stored in a separate collection and linked to posts via PostID.
type Comment struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`      // MongoDB ObjectID
	PostID    primitive.ObjectID `json:"post_id" bson:"post_id"`       // Reference to parent blog post
	Author    string             `json:"author" bson:"author"`         // Comment author name
	Content   string             `json:"content" bson:"content"`       // Comment text content
	CreatedAt time.Time          `json:"created_at" bson:"created_at"` // Creation timestamp
}
