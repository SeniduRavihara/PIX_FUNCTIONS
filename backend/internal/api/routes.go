package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/voltrun/backend/internal/auth"
	"github.com/voltrun/backend/internal/exec"
	"github.com/voltrun/backend/internal/storage"
	"github.com/voltrun/backend/internal/vm"
)

// SetupRoutes registers all API routes
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Auth routes
	authGroup := api.Group("/auth")
	authGroup.Post("/register", handleRegister)
	authGroup.Post("/login", handleLogin)
	authGroup.Post("/refresh", handleRefresh)
	authGroup.Get("/me", auth.AuthRequired(), handleGetCurrentUser)

	// Protected routes (require authentication)
	// Functions routes
	functions := api.Group("/functions")
	functions.Use(auth.AuthRequired())
	functions.Get("/", listFunctions)
	functions.Post("/", createFunction)
	functions.Get("/:id", getFunction)
	functions.Put("/:id", updateFunction)
	functions.Delete("/:id", deleteFunction)
	functions.Post("/:id/execute", executeFunction)

	// Executions routes
	executions := api.Group("/executions")
	executions.Use(auth.AuthRequired())
	executions.Get("/", listExecutions)
	executions.Get("/:id", getExecution)
	executions.Get("/:id/logs", getExecutionLogs)

	// API Keys routes
	keys := api.Group("/keys")
	keys.Use(auth.AuthRequired())
	keys.Get("/", listAPIKeys)
	keys.Post("/", createAPIKey)
	keys.Delete("/:id", deleteAPIKey)
}

// Auth handlers
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func handleRegister(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Check if user already exists
	var existingUser storage.User
	if err := storage.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Email already registered"})
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to process password"})
	}

	// Create user
	user := storage.User{
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
	}

	if err := storage.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.Status(201).JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

func handleLogin(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Find user by email
	var user storage.User
	if err := storage.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// Check password
	if !auth.CheckPassword(req.Password, user.Password) {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

func handleRefresh(c *fiber.Ctx) error {
	userID, err := auth.GetUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var user storage.User
	if err := storage.DB.First(&user, "id = ?", userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Generate new JWT token
	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

func handleGetCurrentUser(c *fiber.Ctx) error {
	userID, err := auth.GetUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var user storage.User
	if err := storage.DB.First(&user, "id = ?", userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

// Function handlers
func listFunctions(c *fiber.Ctx) error {
	userID, err := auth.GetUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var functions []storage.Function
	result := storage.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&functions)
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

	// Get authenticated user ID from JWT token
	userID, err := auth.GetUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
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
		UserID:      userID,
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
	userID, err := auth.GetUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	id := c.Params("id")
	var function storage.Function
	
	if err := storage.DB.Where("id = ? AND user_id = ?", id, userID).First(&function).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Function not found"})
	}
	
	return c.JSON(function)
}

type UpdateFunctionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Code        string `json:"code"`
	EntryPoint  string `json:"entry_point"`
	MemoryMB    int    `json:"memory_mb"`
	TimeoutSec  int    `json:"timeout_sec"`
	Status      string `json:"status"`
}

func updateFunction(c *fiber.Ctx) error {
	userID, err := auth.GetUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	id := c.Params("id")
	var function storage.Function
	
	if err := storage.DB.Where("id = ? AND user_id = ?", id, userID).First(&function).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Function not found"})
	}

	var req UpdateFunctionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Update only provided fields
	if req.Name != "" {
		function.Name = req.Name
	}
	if req.Description != "" {
		function.Description = req.Description
	}
	if req.Code != "" {
		function.Code = req.Code
	}
	if req.EntryPoint != "" {
		function.EntryPoint = req.EntryPoint
	}
	if req.MemoryMB > 0 {
		function.MemoryMB = req.MemoryMB
	}
	if req.TimeoutSec > 0 {
		function.TimeoutSec = req.TimeoutSec
	}
	if req.Status != "" {
		function.Status = req.Status
	}

	if err := storage.DB.Save(&function).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update function"})
	}

	return c.JSON(function)
}

func deleteFunction(c *fiber.Ctx) error {
	userID, err := auth.GetUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	id := c.Params("id")
	var function storage.Function
	
	if err := storage.DB.Where("id = ? AND user_id = ?", id, userID).First(&function).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Function not found"})
	}

	// Delete associated executions first
	storage.DB.Where("function_id = ?", id).Delete(&storage.Execution{})

	// Delete function
	if err := storage.DB.Delete(&function).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete function"})
	}

	return c.JSON(fiber.Map{"message": "Function deleted successfully"})
}

type ExecuteFunctionRequest struct {
	Input map[string]interface{} `json:"input"`
}

func executeFunction(c *fiber.Ctx) error {
	userID, err := auth.GetUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	id := c.Params("id")
	var function storage.Function
	
	if err := storage.DB.Where("id = ? AND user_id = ?", id, userID).First(&function).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Function not found"})
	}

	var req ExecuteFunctionRequest
	if err := c.BodyParser(&req); err != nil {
		// Default to empty input if not provided
		req.Input = make(map[string]interface{})
	}

	// Marshal input to JSON bytes for JSONB
	inputJSON, err := json.Marshal(req.Input)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input format"})
	}

	// Create execution record
	executionID := uuid.New()
	execution := storage.Execution{
		ID:         executionID,
		UserID:     userID,
		FunctionID: function.ID,
		Status:     "pending",
		Input:      inputJSON,
	}

	if err := storage.DB.Create(&execution).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create execution record"})
	}

	// Execute function asynchronously
	go executeAsync(executionID, function, req.Input)

	return c.Status(201).JSON(fiber.Map{
		"execution_id": executionID,
		"status":       "pending",
		"message":      "Function execution started",
	})
}

// Execution handlers
func listExecutions(c *fiber.Ctx) error {
	userID, err := auth.GetUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	functionID := c.Query("function_id")
	query := storage.DB.Where("user_id = ?", userID)

	if functionID != "" {
		query = query.Where("function_id = ?", functionID)
	}

	var executions []storage.Execution
	if err := query.Order("created_at DESC").Preload("Function").Find(&executions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch executions"})
	}

	return c.JSON(executions)
}

func getExecution(c *fiber.Ctx) error {
	userID, err := auth.GetUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	id := c.Params("id")
	var execution storage.Execution

	if err := storage.DB.Where("id = ? AND user_id = ?", id, userID).Preload("Function").First(&execution).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Execution not found"})
	}

	return c.JSON(execution)
}

func getExecutionLogs(c *fiber.Ctx) error {
	userID, err := auth.GetUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	id := c.Params("id")
	var execution storage.Execution

	if err := storage.DB.Where("id = ? AND user_id = ?", id, userID).First(&execution).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Execution not found"})
	}

	return c.JSON(fiber.Map{
		"execution_id": execution.ID,
		"logs":         execution.Logs,
		"status":       execution.Status,
	})
}

// API Key handlers
func listAPIKeys(c *fiber.Ctx) error {
	userID, err := auth.GetUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var keys []storage.APIKey
	if err := storage.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&keys).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch API keys"})
	}

	return c.JSON(keys)
}

type CreateAPIKeyRequest struct {
	Name string `json:"name" validate:"required"`
}

func createAPIKey(c *fiber.Ctx) error {
	userID, err := auth.GetUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req CreateAPIKeyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Generate random API key
	keyBytes := make([]byte, 32)
	if _, err := rand.Read(keyBytes); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate key"})
	}
	rawKey := "vr_" + hex.EncodeToString(keyBytes)

	// Hash the key for storage
	hashedKey, err := auth.HashPassword(rawKey)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to process key"})
	}

	// Store only prefix for display
	prefix := rawKey[:12]

	apiKey := storage.APIKey{
		UserID: userID,
		Name:   req.Name,
		Key:    hashedKey,
		Prefix: prefix,
	}

	if err := storage.DB.Create(&apiKey).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create API key"})
	}

	// Return the raw key only once
	return c.Status(201).JSON(fiber.Map{
		"id":     apiKey.ID,
		"name":   apiKey.Name,
		"key":    rawKey,
		"prefix": prefix,
		"message": "Save this key securely. It won't be shown again.",
	})
}

func deleteAPIKey(c *fiber.Ctx) error {
	userID, err := auth.GetUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	id := c.Params("id")
	var apiKey storage.APIKey

	if err := storage.DB.Where("id = ? AND user_id = ?", id, userID).First(&apiKey).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "API key not found"})
	}

	if err := storage.DB.Delete(&apiKey).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete API key"})
	}

	return c.JSON(fiber.Map{"message": "API key deleted successfully"})
}

// marshalJSON converts a map to JSON bytes for JSONB
func marshalJSON(data interface{}) []byte {
	bytes, _ := json.Marshal(data)
	return bytes
}

// executeAsync executes a function asynchronously
func executeAsync(executionID uuid.UUID, function storage.Function, input map[string]interface{}) {
	ctx := context.Background()
	
	// Create VM manager and execution engine
	vmManager := vm.NewVMManager()
	engine := exec.NewExecutionEngine(vmManager)
	
	// Execute the function
	result, err := engine.Execute(ctx, exec.ExecutionRequest{
		FunctionID: function.ID,
		Input:      input,
		UserID:     function.UserID,
	})
	
	if err != nil {
		// Update execution with error
		storage.DB.Model(&storage.Execution{}).Where("id = ?", executionID).Updates(map[string]interface{}{
			"status": "failed",
			"error":  err.Error(),
			"logs":   fmt.Sprintf("Execution failed: %v", err),
		})
		return
	}
	
	// Update execution with results
	storage.DB.Model(&storage.Execution{}).Where("id = ?", executionID).Updates(map[string]interface{}{
		"status":      result.Status,
		"output":      marshalJSON(result.Output),
		"logs":        result.Logs,
		"duration_ms": result.DurationMS,
	})
}
