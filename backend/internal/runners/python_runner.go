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

// PythonRunner executes Python functions
type PythonRunner struct{}

// Execute runs Python code with the given input
func (r *PythonRunner) Execute(ctx context.Context, code string, input map[string]interface{}, timeout time.Duration) (*ExecutionResult, error) {
	start := time.Now()

	// Create temporary directory for execution
	tempDir, err := os.MkdirTemp("", "voltrun-python-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Write code to file
	codePath := filepath.Join(tempDir, "handler.py")
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
import json
import sys
import traceback
from handler import handler

if __name__ == '__main__':
    try:
        with open('input.json', 'r') as f:
            event = json.load(f)
        
        result = handler(event)
        
        print('__VOLTRUN_OUTPUT_START__')
        print(json.dumps(result))
        print('__VOLTRUN_OUTPUT_END__')
    except Exception as e:
        print(f'__VOLTRUN_ERROR__ {str(e)}', file=sys.stderr)
        traceback.print_exc(file=sys.stderr)
        sys.exit(1)
`

	wrapperPath := filepath.Join(tempDir, "wrapper.py")
	if err := os.WriteFile(wrapperPath, []byte(wrapperCode), 0644); err != nil {
		return nil, fmt.Errorf("failed to write wrapper: %w", err)
	}

	// Execute Python with timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctxWithTimeout, "python3", "wrapper.py")
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
