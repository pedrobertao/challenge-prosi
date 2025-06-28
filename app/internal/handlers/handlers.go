// Package handlers provides HTTP request handlers for the blog API.
// It contains all the endpoint handlers for managing blog posts and comments,
// including CRUD operations and proper error handling with JSON responses.
package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pedrobertao/challenge-prosi/app/internal/models"
	"github.com/pedrobertao/challenge-prosi/app/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Handler struct holds the database storage instance and provides
// methods for handling HTTP requests to the blog API endpoints.
type Handler struct {
	DB *storage.Storage // Database storage instance for MongoDB operations
}

// New creates and returns a new Handler instance with the provided storage.
// This is the constructor function for the Handler struct.
//
// Parameters:
//   - db: pointer to a Storage instance for database operations
//
// Returns a pointer to a new Handler instance.
func New(db *storage.Storage) *Handler {
	return &Handler{DB: db}
}

// GetPosts handles GET /api/posts requests.
// Returns a list of all blog posts with summary information including
// post ID, title, comment count, and creation date.
// This endpoint provides an overview of all posts without full content.
func (h *Handler) GetPosts(c *fiber.Ctx) error {
	// Create context with timeout to prevent hanging database operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch all posts from the database
	cursor, err := h.DB.Posts.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(500).JSON(models.APIResponse{
			Success: false,
			Error:   "Failed to fetch posts",
		})
	}
	defer cursor.Close(ctx)

	// Build summary list with comment counts for each post
	var summaries []models.BlogPostSummary
	for cursor.Next(ctx) {
		var post models.BlogPost
		if err := cursor.Decode(&post); err != nil {
			// Skip malformed posts and continue processing
			continue
		}

		// Count comments for this post
		count, _ := h.DB.Comments.CountDocuments(ctx, bson.M{"post_id": post.ID})
		summaries = append(summaries, models.BlogPostSummary{
			ID:           post.ID,
			Title:        post.Title,
			CommentCount: count,
			CreatedAt:    post.CreatedAt,
		})
	}

	return c.JSON(models.APIResponse{Success: true, Data: summaries})
}

// CreatePost handles POST /api/posts requests.
// Creates a new blog post with the provided title and content.
// Validates required fields and returns the created post with its generated ID.
func (h *Handler) CreatePost(c *fiber.Ctx) error {
	// Parse the request body into the expected structure
	var req models.CreatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid JSON",
		})
	}

	// Validate required fields
	if req.Title == "" || req.Content == "" {
		return c.Status(400).JSON(models.APIResponse{
			Success: false,
			Error:   "Title and content required",
		})
	}

	// Create context with timeout for database operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create new blog post with current timestamp
	post := models.BlogPost{
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	// Insert the post into the database
	result, err := h.DB.Posts.InsertOne(ctx, post)
	if err != nil {
		return c.Status(500).JSON(models.APIResponse{
			Success: false,
			Error:   "Failed to create post",
		})
	}

	// Set the generated ID and return the complete post
	post.ID = result.InsertedID.(primitive.ObjectID)
	return c.JSON(models.APIResponse{Success: true, Data: post})
}

// GetPost handles GET /api/posts/:id requests.
// Retrieves a specific blog post by its ID along with all associated comments.
// Returns 404 if the post doesn't exist, or 400 if the ID format is invalid.
func (h *Handler) GetPost(c *fiber.Ctx) error {
	// Parse and validate the post ID from URL parameters
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid post ID",
		})
	}

	// Create context with timeout for database operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find the specific post by ID
	var post models.BlogPost
	err = h.DB.Posts.FindOne(ctx, bson.M{"_id": id}).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(models.APIResponse{
				Success: false,
				Error:   "Post not found",
			})
		}
		return c.Status(500).JSON(models.APIResponse{
			Success: false,
			Error:   "Failed to fetch post",
		})
	}

	// Fetch all comments for this post and attach them
	cursor, err := h.DB.Comments.Find(ctx, bson.M{"post_id": id})
	if err == nil {
		cursor.All(ctx, &post.Comments)
		cursor.Close(ctx)
	}

	return c.JSON(models.APIResponse{Success: true, Data: post})
}

// CreateComment handles POST /api/posts/:id/comments requests.
// Creates a new comment on a specific blog post.
// Validates that the post exists and that required comment fields are provided.
func (h *Handler) CreateComment(c *fiber.Ctx) error {
	// Parse and validate the post ID from URL parameters
	postID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid post ID",
		})
	}

	// Parse the request body into the expected structure
	var req models.CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid JSON",
		})
	}

	// Validate required comment fields
	if req.Author == "" || req.Content == "" {
		return c.Status(400).JSON(models.APIResponse{
			Success: false,
			Error:   "Author and content required",
		})
	}

	// Create context with timeout for database operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Verify that the target post exists
	count, err := h.DB.Posts.CountDocuments(ctx, bson.M{"_id": postID})
	if err != nil || count == 0 {
		return c.Status(404).JSON(models.APIResponse{
			Success: false,
			Error:   "Post not found",
		})
	}

	// Create new comment with current timestamp
	comment := models.Comment{
		PostID:    postID,
		Author:    req.Author,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	// Insert the comment into the database
	result, err := h.DB.Comments.InsertOne(ctx, comment)
	if err != nil {
		return c.Status(500).JSON(models.APIResponse{
			Success: false,
			Error:   "Failed to create comment",
		})
	}

	// Set the generated ID and return the complete comment
	comment.ID = result.InsertedID.(primitive.ObjectID)
	return c.JSON(models.APIResponse{Success: true, Data: comment})
}
