package routes

import (
	"class-go-ai/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App) {
	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "hello world",
		})
	})

	// User routes
	app.Get("/users", handlers.GetUsers)
	app.Get("/users/:id", handlers.GetUser)
	app.Post("/users", handlers.CreateUser)
	app.Put("/users/:id", handlers.UpdateUser)
	app.Delete("/users/:id", handlers.DeleteUser)

	// Transfer routes
	app.Post("/transfers", handlers.CreateTransfer)
	app.Get("/transfers/:id", handlers.GetTransfer)
	app.Get("/transfers", handlers.ListTransfers)
}
