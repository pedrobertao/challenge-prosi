package http

import (
	"github.com/gofiber/fiber/v2"
)

func Start() error {
	// Initialize a new Fiber app
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Test App v1.0.1",
	})

	// Define a route for the GET method on the root path '/'
	app.Get("/", func(c *fiber.Ctx) error {
		return fiber.NewError(782, "Custom error message")
	})
	// S

	return app.Listen(":3030")
}
