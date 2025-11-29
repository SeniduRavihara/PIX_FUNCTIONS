package runners

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// NodeRunner executes Node.js functions
type NodeRunner struct{}

// ExecutionResult represents the result of code execution
type ExecutionResult struct {
	Output     map[string]interface{} `json:"output"`
	Logs       string                 `json:"logs"`
	Error      string                 `json:"error,omitempty"`
	DurationMS int64                  `json:"duration_ms"`
	ExitCode   int                    `json:"exit_code"`
}

// Execute runs Node.js code with the given input
func (r *NodeRunner) Execute(ctx context.Context, code string, input map[string]interface{}, timeout time.Duration) (*ExecutionResult, error) {
	start := time.Now()

	// Create temporary directory for execution
	tempDir, err := os.MkdirTemp("", "voltrun-node-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Write code to file
	codePath := filepath.Join(tempDir, "index.js")
	if err := os.WriteFile(codePath, []byte(code), 0644); err != nil {
		return nil, fmt.Errorf("failed to write code: %w", err)
	}

	// Write input to file
	inputJSON, _ := json.Marshal(input)
	inputPath := filepath.Join(tempDir, "input.json")
	if err := os.WriteFile(inputPath, inputJSON, 0644); err != nil {
		return nil, fmt.Errorf("failed to write input: %w", err)
	}

	// Create wrapper script that loads input and executes handler
	wrapperCode := `
const fs = require('fs');
const handler = require('./index.js');

const input = JSON.parse(fs.readFileSync('./input.json', 'utf8'));

(async () => {
  try {
    const result = await handler.handler(input);
    console.log('__VOLTRUN_OUTPUT_START__');
    console.log(JSON.stringify(result));
    console.log('__VOLTRUN_OUTPUT_END__');
  } catch (error) {
    console.error('__VOLTRUN_ERROR__', error.message);
    console.error(error.stack);
    process.exit(1);
  }
})();
`

	wrapperPath := filepath.Join(tempDir, "wrapper.js")
	if err := os.WriteFile(wrapperPath, []byte(wrapperCode), 0644); err != nil {
		return nil, fmt.Errorf("failed to write wrapper: %w", err)
	}

	// Execute Node.js with timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctxWithTimeout, "node", "wrapper.js")
	cmd.Dir = tempDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	duration := time.Since(start).Milliseconds()

	result := &ExecutionResult{
		DurationMS: duration,
		Logs:       stdout.String() + stderr.String(),
	}

	if err != nil {
		if ctxWithTimeout.Err() == context.DeadlineExceeded {
			result.Error = "Execution timeout exceeded"
			result.ExitCode = -1
		} else {
			result.Error = err.Error()
			if exitErr, ok := err.(*exec.ExitError); ok {
				result.ExitCode = exitErr.ExitCode()
			}
		}
		return result, nil
	}

	// Parse output from stdout
	outputStr := stdout.String()
	if startIdx := bytes.Index([]byte(outputStr), []byte("__VOLTRUN_OUTPUT_START__")); startIdx != -1 {
		if endIdx := bytes.Index([]byte(outputStr), []byte("__VOLTRUN_OUTPUT_END__")); endIdx != -1 {
			jsonStart := startIdx + len("__VOLTRUN_OUTPUT_START__") + 1
			outputJSON := outputStr[jsonStart:endIdx]
			var output map[string]interface{}
			if err := json.Unmarshal([]byte(outputJSON), &output); err == nil {
				result.Output = output
			} else {
				result.Output = map[string]interface{}{"result": outputJSON}
			}
		}
	}

	if result.Output == nil {
		result.Output = map[string]interface{}{"message": "Execution completed"}
	}

	result.ExitCode = 0
	return result, nil
}
