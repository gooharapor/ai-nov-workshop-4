package main

import (
	"log"

	"class-go-ai/database"
	"class-go-ai/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Initialize database connection
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Create new Fiber app
	app := fiber.New(fiber.Config{
		AppName: "User Management API v1.0",
	})

	// Middlewares
	app.Use(logger.New())
	app.Use(cors.New())

	// Setup routes
	routes.SetupRoutes(app)

	// Start server
	log.Println("Server starting on port 3000...")
	log.Fatal(app.Listen(":3000"))
}
