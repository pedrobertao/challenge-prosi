package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pedrobertao/challenge-prosi/app/internal/handlers"
)

func Setup(h *handlers.Handler) *fiber.App {
	app := fiber.New()
	registerRoutes(app, h)

	return app
}

func registerRoutes(app *fiber.App, h *handlers.Handler) fiber.Router {
	router := app.Group("/api")

	app.Get("/api/posts", h.GetPosts)
	app.Get("/api/posts/:id", h.GetPost)
	app.Post("/api/posts", h.CreatePost)
	app.Post("/api/posts/:id/comments", h.CreateComment)

	return router
}
