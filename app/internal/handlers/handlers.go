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

type Handler struct {
	DB *storage.Storage
}

func New(db *storage.Storage) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) GetPosts(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := h.DB.Posts.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(500).JSON(models.APIResponse{
			Success: false,
			Error:   "Failed to fetch posts",
		})
	}
	defer cursor.Close(ctx)

	var summaries []models.BlogPostSummary
	for cursor.Next(ctx) {
		var post models.BlogPost
		if err := cursor.Decode(&post); err != nil {
			continue
		}

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

// POST /api/posts
func (h *Handler) CreatePost(c *fiber.Ctx) error {
	var req models.CreatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid JSON",
		})
	}

	if req.Title == "" || req.Content == "" {
		return c.Status(400).JSON(models.APIResponse{
			Success: false,
			Error:   "Title and content required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	post := models.BlogPost{
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	result, err := h.DB.Posts.InsertOne(ctx, post)
	if err != nil {
		return c.Status(500).JSON(models.APIResponse{
			Success: false,
			Error:   "Failed to create post",
		})
	}

	post.ID = result.InsertedID.(primitive.ObjectID)
	return c.JSON(models.APIResponse{Success: true, Data: post})
}

// GET /api/posts/:id
func (h *Handler) GetPost(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid post ID",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

	cursor, err := h.DB.Comments.Find(ctx, bson.M{"post_id": id})
	if err == nil {
		cursor.All(ctx, &post.Comments)
		cursor.Close(ctx)
	}

	return c.JSON(models.APIResponse{Success: true, Data: post})
}

// POST /api/posts/:id/comments
func (h *Handler) CreateComment(c *fiber.Ctx) error {
	postID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid post ID",
		})
	}

	var req models.CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid JSON",
		})
	}

	if req.Author == "" || req.Content == "" {
		return c.Status(400).JSON(models.APIResponse{
			Success: false,
			Error:   "Author and content required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := h.DB.Posts.CountDocuments(ctx, bson.M{"_id": postID})
	if err != nil || count == 0 {
		return c.Status(404).JSON(models.APIResponse{
			Success: false,
			Error:   "Post not found",
		})
	}

	comment := models.Comment{
		PostID:    postID,
		Author:    req.Author,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	result, err := h.DB.Comments.InsertOne(ctx, comment)
	if err != nil {
		return c.Status(500).JSON(models.APIResponse{
			Success: false,
			Error:   "Failed to create comment",
		})
	}

	comment.ID = result.InsertedID.(primitive.ObjectID)
	return c.JSON(models.APIResponse{Success: true, Data: comment})
}
