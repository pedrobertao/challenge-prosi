// Package routes handles HTTP route configuration and setup for the blog API.
// It provides functions to create and configure the Fiber application with
// all necessary endpoints for blog posts and comments management.
package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pedrobertao/challenge-prosi/app/internal/handlers"
)

// Setup creates and configures a new Fiber application with all API routes.
// This is the main entry point for setting up the HTTP server with proper
// route configuration and handler registration.
//
// Parameters:
//   - h: pointer to a Handler instance containing all endpoint handlers
//
// Returns a configured Fiber application ready to serve HTTP requests.
func Setup(h *handlers.Handler) *fiber.App {
	// Create a new Fiber application instance with default configuration
	app := fiber.New()

	// Register all API routes with their corresponding handlers
	registerRoutes(app, h)

	return app
}

// registerRoutes configures all API endpoints for the blog application.
// Sets up RESTful routes for blog posts and comments under the /api prefix.
// This function organizes all route definitions in one place for maintainability.
//
// API Endpoints configured:
//   - GET    /api/posts           - List all blog posts (summary view)
//   - GET    /api/posts/:id       - Get specific post with comments
//   - POST   /api/posts           - Create a new blog post
//   - POST   /api/posts/:id/comments - Add comment to a specific post
//
// Parameters:
//   - app: the Fiber application instance to register routes on
//   - h: pointer to Handler instance containing endpoint implementations
//
// Returns the API router group for potential additional configuration.
func registerRoutes(app *fiber.App, h *handlers.Handler) fiber.Router {
	// Create API route group for all endpoints under /api prefix
	router := app.Group("/api")

	// Blog posts endpoints
	app.Get("/api/posts", h.GetPosts)    // List all posts with summaries
	app.Get("/api/posts/:id", h.GetPost) // Get single post with comments
	app.Post("/api/posts", h.CreatePost) // Create new blog post

	// Comments endpoint
	app.Post("/api/posts/:id/comments", h.CreateComment) // Add comment to post

	return router
}
