package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"github.com/voltrun/backend/internal/api"
	"github.com/voltrun/backend/internal/storage"
	"github.com/voltrun/backend/internal/utils"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Load configuration
	config := utils.LoadConfig()

	// Initialize logger
	if err := utils.InitLogger(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer utils.Logger.Sync()

	utils.Info("Starting VoltRun backend server")

	// Initialize database
	if err := storage.InitDB(); err != nil {
		utils.Error("Failed to initialize database")
		log.Fatalf("Database initialization failed: %v", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "VoltRun v1.0.0",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			utils.Error("Request error")
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "voltrun-backend",
			"version": "1.0.0",
		})
	})

	// Setup API routes
	api.SetupRoutes(app)

	// Start server
	utils.Info("Server starting on port " + config.Port)
	log.Printf("🚀 VoltRun server listening on http://localhost:%s", config.Port)
	log.Printf("📊 Health check: http://localhost:%s/health", config.Port)
	log.Printf("� API endpoint: http://localhost:%s/api", config.Port)
	
	if err := app.Listen(":" + config.Port); err != nil {
		utils.Error("Server failed to start")
		log.Fatal(err)
	}
}
