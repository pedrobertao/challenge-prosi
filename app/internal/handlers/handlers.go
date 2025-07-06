// Package handlers provides HTTP request handlers for the blog API.
// It contains all the endpoint handlers for managing blog posts and comments,
// including CRUD operations and proper error handling with JSON responses.
package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pedrobertao/challenge-prosi/app/internal/models"
	"github.com/pedrobertao/challenge-prosi/app/internal/storage"
	"github.com/pedrobertao/challenge-prosi/app/lib/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// DEFAULT_DB_TIMEOUT defines the maximum duration for database operations
// to prevent hanging requests and ensure responsive API behavior.
const DEFAULT_DB_TIMEOUT = 10 * time.Second

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
//
// Response format:
//   - 200: Success with array of BlogPostSummary objects
//   - 404: No posts found (returns empty array)
//   - 502: Database connection or query error
func (h *Handler) GetPosts(c *fiber.Ctx) error {
	// Create context with timeout to prevent hanging database operations
	ctx, cancel := context.WithTimeout(c.Context(), DEFAULT_DB_TIMEOUT)
	defer cancel()

	// Fetch all posts from the database
	cursor, err := h.DB.Posts.Find(ctx, bson.M{})
	if err != nil {
		if err == mongo.ErrEmptySlice {
			return c.Status(http.StatusNotFound).JSON(models.APIResponse{
				Success: true,
				Data:    []models.BlogPostSummary{},
				Error:   "",
			})
		}
		return c.Status(http.StatusBadGateway).JSON(models.APIResponse{
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
			logger.Warn("malformed post", zap.Error(err))
			// Skip malformed posts and continue processing
			continue
		}

		// Count comments for this post
		count, err := h.DB.Comments.CountDocuments(ctx, bson.M{"post_id": post.ID})
		if err != nil {
			logger.Error("failed to count objects", zap.Error(err))
		}
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
//
// Request body should contain:
//   - title: string (required) - The post title
//   - content: string (required) - The post content
//
// Response format:
//   - 200: Success with created BlogPost object
//   - 400: Invalid JSON or missing required fields
//   - 502: Database insertion error
func (h *Handler) CreatePost(c *fiber.Ctx) error {
	// Parse the request body into the expected structure
	var req models.CreatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
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
	ctx, cancel := context.WithTimeout(c.Context(), DEFAULT_DB_TIMEOUT)
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
		return c.Status(http.StatusBadGateway).JSON(models.APIResponse{
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
//
// URL parameters:
//   - id: string (required) - MongoDB ObjectID as hex string
//
// Response format:
//   - 200: Success with BlogPost object including comments array
//   - 400: Invalid ObjectID format
//   - 404: Post not found
//   - 500: Database query error
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
	ctx, cancel := context.WithTimeout(c.Context(), DEFAULT_DB_TIMEOUT)
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

// DeletePost handles DELETE /api/posts/:id requests.
// Deletes a specific blog post and all its associated comments atomically
// using MongoDB transactions to ensure data consistency.
//
// URL parameters:
//   - id: string (required) - MongoDB ObjectID as hex string
//
// Response format:
//   - 200: Success - post and comments deleted
//   - 400: Invalid ObjectID format or post not found
//   - 502: Database transaction or deletion error
//
// This operation uses MongoDB transactions to ensure that both the post
// and all its comments are deleted together, preventing orphaned comments.
func (h *Handler) DeletePost(c *fiber.Ctx) error {
	// Parse and validate the post ID from URL parameters
	postID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid post ID",
		})
	}

	// Create context with timeout for database operations
	ctx, cancel := context.WithTimeout(c.Context(), DEFAULT_DB_TIMEOUT)
	defer cancel()

	// Start a session for transaction to ensure atomicity
	session, err := h.DB.Client.StartSession()
	if err != nil {
		logger.Error("failed to start session from db", zap.Error(err))
		return c.Status(http.StatusBadGateway).JSON(models.APIResponse{
			Success: false,
			Error:   "Failed to delete post",
		})
	}
	defer session.EndSession(ctx)

	// Initialize response variables for transaction error handling
	status := http.StatusOK
	response := models.APIResponse{
		Data:    nil,
		Success: false,
		Error:   "",
	}

	// Execute transaction - both operations must succeed or both will rollback
	if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		// Step 1: Delete all comments associated with this post
		commentFilter := bson.M{"post_id": postID}
		_, err := h.DB.Comments.DeleteMany(sc, commentFilter)
		if err != nil {
			logger.Error("failed to delete comments from session", zap.Error(err))
			status = http.StatusBadGateway
			response.Error = "Failed to delete comments from post"
			return err
		}

		// Step 2: Delete the blog post itself
		postFilter := bson.M{"_id": postID}
		deletePostResult, err := h.DB.Posts.DeleteOne(sc, postFilter)
		if err != nil {
			logger.Error("failed to delete post from session", zap.Error(err))
			status = http.StatusBadGateway
			response.Error = "Failed to delete post"
			return err
		}

		// Verify that the post actually existed and was deleted
		if deletePostResult.DeletedCount == 0 {
			status = http.StatusBadRequest
			response.Error = "Post not found"
			return errors.New("no post deleted")
		}
		return nil
	}); err != nil {
		// Transaction failed - return the error details
		return c.Status(status).JSON(models.APIResponse{
			Success: false,
			Error:   response.Error,
		})
	}

	// Transaction succeeded - post and comments deleted
	return c.Status(status).JSON(models.APIResponse{Data: postID, Success: true, Error: ""})
}

// CreateComment handles POST /api/posts/:id/comments requests.
// Creates a new comment on a specific blog post.
// Validates that the post exists and that required comment fields are provided.
//
// URL parameters:
//   - id: string (required) - MongoDB ObjectID of the target post
//
// Request body should contain:
//   - author: string (required) - Comment author name
//   - content: string (required) - Comment content
//
// Response format:
//   - 200: Success with created Comment object
//   - 400: Invalid JSON, missing fields, or invalid post ID
//   - 404: Target post not found
//   - 500: Database insertion error
func (h *Handler) CreateComment(c *fiber.Ctx) error {
	// Parse and validate the post ID from URL parameters
	postID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid post ID",
		})
	}

	// Parse the request body into the expected structure
	var req models.CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
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
	ctx, cancel := context.WithTimeout(c.Context(), DEFAULT_DB_TIMEOUT)
	defer cancel()

	// Verify that the target post exists before creating comment
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

// DeleteComment handles DELETE /api/comments/:id requests.
// Deletes a specific comment by its ID.
//
// URL parameters:
//   - id: string (required) - MongoDB ObjectID of the comment to delete
//
// Response format:
//   - 200: Success - comment deleted
//   - 400: Invalid ObjectID format or comment not found
//   - 502: Database deletion error
//
// Note: This operation only deletes the comment itself and does not
// require transaction handling since it's a single atomic operation.
func (h *Handler) DeleteComment(c *fiber.Ctx) error {
	// Parse and validate the comment ID from URL parameters
	commentID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid comment ID",
		})
	}

	// Create context with timeout for database operations
	ctx, cancel := context.WithTimeout(c.Context(), DEFAULT_DB_TIMEOUT)
	defer cancel()

	// Create filter using the comment ID for deletion
	filter := bson.M{"_id": commentID}

	// Execute the deletion operation
	result, err := h.DB.Comments.DeleteOne(ctx, filter)
	if err != nil {
		logger.Error("failed to delete comment", zap.Error(err))
		return c.Status(http.StatusBadGateway).JSON(models.APIResponse{
			Success: false,
			Error:   "Failed to delete comment",
		})
	}

	// Check if a comment was actually found and deleted
	if result.DeletedCount == 0 {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "No comment found to delete",
		})
	}

	// Successfully deleted the comment
	return c.Status(http.StatusOK).JSON(models.APIResponse{
		Data:    commentID,
		Success: true,
		Error:   "",
	})
}
