package exec

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/voltrun/backend/internal/runners"
	"github.com/voltrun/backend/internal/storage"
	"github.com/voltrun/backend/internal/vm"
)

// ExecutionEngine handles function execution
type ExecutionEngine struct {
	vmManager *vm.VMManager
}

// NewExecutionEngine creates a new execution engine
func NewExecutionEngine(vmManager *vm.VMManager) *ExecutionEngine {
	return &ExecutionEngine{
		vmManager: vmManager,
	}
}

// ExecutionRequest represents a function execution request
type ExecutionRequest struct {
	FunctionID uuid.UUID              `json:"function_id"`
	Input      map[string]interface{} `json:"input"`
	UserID     uuid.UUID              `json:"user_id"`
}

// ExecutionResult represents the result of a function execution
type ExecutionResult struct {
	ExecutionID uuid.UUID              `json:"execution_id"`
	Output      map[string]interface{} `json:"output"`
	Logs        string                 `json:"logs"`
	Error       string                 `json:"error,omitempty"`
	DurationMS  int64                  `json:"duration_ms"`
	Status      string                 `json:"status"`
}

// Execute runs a function in an isolated VM
func (e *ExecutionEngine) Execute(ctx context.Context, req ExecutionRequest) (*ExecutionResult, error) {
	startTime := time.Now()

	// Fetch function from database
	function, err := e.getFunction(req.FunctionID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch function: %w", err)
	}

	// Create execution record
	execution := &storage.Execution{
		ID:         uuid.New(),
		UserID:     req.UserID,
		FunctionID: req.FunctionID,
		Status:     "pending",
		Input:      marshalJSON(req.Input),
		CreatedAt:  time.Now(),
	}

	if err := storage.DB.Create(execution).Error; err != nil {
		return nil, fmt.Errorf("failed to create execution record: %w", err)
	}

	// Update status to running
	now := time.Now()
	execution.Status = "running"
	execution.StartedAt = &now
	storage.DB.Save(execution)

	// Create VM configuration
	vmConfig := vm.VMConfig{
		ID:          vm.GenerateVMID(),
		MemoryMB:    function.MemoryMB,
		CPUs:        1,
		TimeoutSec:  function.TimeoutSec,
		Environment: map[string]string{},
	}

	// Create and start VM
	vmInstance, err := e.vmManager.CreateVM(ctx, vmConfig)
	if err != nil {
		e.updateExecutionError(execution, fmt.Sprintf("VM creation failed: %v", err))
		return nil, fmt.Errorf("failed to create VM: %w", err)
	}

	// Execute function inside VM
	result, err := e.executeInVM(ctx, vmInstance, function, req.Input)
	
	// Cleanup: destroy VM
	defer e.vmManager.DestroyVM(ctx, vmInstance.ID)

	duration := time.Since(startTime).Milliseconds()
	completedAt := time.Now()

	if err != nil {
		execution.Status = "failed"
		execution.Error = err.Error()
		execution.DurationMS = duration
		execution.CompletedAt = &completedAt
		storage.DB.Save(execution)

		return &ExecutionResult{
			ExecutionID: execution.ID,
			Status:      "failed",
			Error:       err.Error(),
			DurationMS:  duration,
		}, nil
	}

	// Update execution record with results
	execution.Status = "success"
	execution.Output = marshalJSON(result.Output)
	execution.Logs = result.Logs
	execution.DurationMS = duration
	execution.CompletedAt = &completedAt
	storage.DB.Save(execution)

	return &ExecutionResult{
		ExecutionID: execution.ID,
		Output:      result.Output,
		Logs:        result.Logs,
		DurationMS:  duration,
		Status:      "success",
	}, nil
}

// executeInVM executes code inside a VM (placeholder)
func (e *ExecutionEngine) executeInVM(ctx context.Context, vm *vm.VM, function *storage.Function, input map[string]interface{}) (*ExecutionResult, error) {
	// TODO: Implement actual code execution inside Firecracker VM
	// This is a placeholder that simulates execution
	
	// Based on runtime, dispatch to appropriate runner
	switch function.Runtime {
	case "nodejs":
		return e.executeNodeJS(ctx, function, input)
	case "python":
		return e.executePython(ctx, function, input)
	default:
		return nil, fmt.Errorf("unsupported runtime: %s", function.Runtime)
	}
}

// executeNodeJS executes Node.js code
func (e *ExecutionEngine) executeNodeJS(ctx context.Context, function *storage.Function, input map[string]interface{}) (*ExecutionResult, error) {
	runner := &runners.NodeRunner{}
	timeout := time.Duration(function.TimeoutSec) * time.Second
	
	result, err := runner.Execute(ctx, function.Code, input, timeout)
	if err != nil {
		return nil, fmt.Errorf("node execution failed: %w", err)
	}
	
	return &ExecutionResult{
		Output:     result.Output,
		Logs:       result.Logs,
		DurationMS: result.DurationMS,
	}, nil
}

// executePython executes Python code
func (e *ExecutionEngine) executePython(ctx context.Context, function *storage.Function, input map[string]interface{}) (*ExecutionResult, error) {
	runner := &runners.PythonRunner{}
	timeout := time.Duration(function.TimeoutSec) * time.Second
	
	result, err := runner.Execute(ctx, function.Code, input, timeout)
	if err != nil {
		return nil, fmt.Errorf("python execution failed: %w", err)
	}
	
	return &ExecutionResult{
		Output:     result.Output,
		Logs:       result.Logs,
		DurationMS: result.DurationMS,
	}, nil
}

// getFunction retrieves a function from the database
func (e *ExecutionEngine) getFunction(functionID uuid.UUID) (*storage.Function, error) {
	var function storage.Function
	if err := storage.DB.First(&function, "id = ?", functionID).Error; err != nil {
		return nil, err
	}
	return &function, nil
}

// updateExecutionError updates an execution with error status
func (e *ExecutionEngine) updateExecutionError(execution *storage.Execution, errorMsg string) {
	now := time.Now()
	execution.Status = "failed"
	execution.Error = errorMsg
	execution.CompletedAt = &now
	storage.DB.Save(execution)
}

// marshalJSON converts a map to JSON string
func marshalJSON(data interface{}) string {
	bytes, _ := json.Marshal(data)
	return string(bytes)
}
