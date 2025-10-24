package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/voltrun/backend/internal/storage"
)

// SetupRoutes registers all API routes
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", handleRegister)
	auth.Post("/login", handleLogin)
	auth.Post("/refresh", handleRefresh)

	// Protected routes (require authentication)
	// Functions routes
	functions := api.Group("/functions")
	// functions.Use(middleware.AuthRequired())
	functions.Get("/", listFunctions)
	functions.Post("/", createFunction)
	functions.Get("/:id", getFunction)
	functions.Put("/:id", updateFunction)
	functions.Delete("/:id", deleteFunction)
	functions.Post("/:id/execute", executeFunction)

	// Executions routes
	executions := api.Group("/executions")
	// executions.Use(middleware.AuthRequired())
	executions.Get("/", listExecutions)
	executions.Get("/:id", getExecution)
	executions.Get("/:id/logs", getExecutionLogs)

	// API Keys routes
	keys := api.Group("/keys")
	// keys.Use(middleware.AuthRequired())
	keys.Get("/", listAPIKeys)
	keys.Post("/", createAPIKey)
	keys.Delete("/:id", deleteAPIKey)
}

// Auth handlers
func handleRegister(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "register endpoint - to be implemented"})
}

func handleLogin(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "login endpoint - to be implemented"})
}

func handleRefresh(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "refresh endpoint - to be implemented"})
}

// Function handlers
func listFunctions(c *fiber.Ctx) error {
	var functions []storage.Function
	result := storage.DB.Order("created_at DESC").Find(&functions)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch functions"})
	}
	return c.JSON(functions)
}

type CreateFunctionRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Runtime     string `json:"runtime" validate:"required"`
	Code        string `json:"code" validate:"required"`
	EntryPoint  string `json:"entry_point"`
	MemoryMB    int    `json:"memory_mb"`
	TimeoutSec  int    `json:"timeout_sec"`
}

func createFunction(c *fiber.Ctx) error {
	var req CreateFunctionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// TODO: Get actual user ID from JWT token
	// For now, create a demo user if not exists
	var user storage.User
	result := storage.DB.Where("email = ?", "demo@voltrun.com").First(&user)
	if result.Error != nil {
		// Create demo user
		user = storage.User{
			Email:    "demo@voltrun.com",
			Password: "demo", // In production, this should be hashed
			Name:     "Demo User",
		}
		if err := storage.DB.Create(&user).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
		}
	}

	// Set defaults
	if req.EntryPoint == "" {
		req.EntryPoint = "index.handler"
	}
	if req.MemoryMB == 0 {
		req.MemoryMB = 128
	}
	if req.TimeoutSec == 0 {
		req.TimeoutSec = 30
	}

	function := storage.Function{
		UserID:      user.ID,
		Name:        req.Name,
		Description: req.Description,
		Runtime:     req.Runtime,
		Code:        req.Code,
		EntryPoint:  req.EntryPoint,
		MemoryMB:    req.MemoryMB,
		TimeoutSec:  req.TimeoutSec,
		Status:      "active",
	}

	if err := storage.DB.Create(&function).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create function"})
	}

	return c.Status(201).JSON(function)
}

func getFunction(c *fiber.Ctx) error {
	id := c.Params("id")
	var function storage.Function
	
	if err := storage.DB.First(&function, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Function not found"})
	}
	
	return c.JSON(function)
}

func updateFunction(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "update function - to be implemented"})
}

func deleteFunction(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "delete function - to be implemented"})
}

func executeFunction(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "execute function - to be implemented"})
}

// Execution handlers
func listExecutions(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "list executions - to be implemented"})
}

func getExecution(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "get execution - to be implemented"})
}

func getExecutionLogs(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "get execution logs - to be implemented"})
}

// API Key handlers
func listAPIKeys(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "list API keys - to be implemented"})
}

func createAPIKey(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "create API key - to be implemented"})
}

func deleteAPIKey(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "delete API key - to be implemented"})
}
