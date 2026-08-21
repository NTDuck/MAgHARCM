package tools

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/cloudwego/eino/adk/filesystem"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"

	. "MAgHARCM/internal/patterns"
)

// ExecuteInput parameters for executing shell commands.
type ExecuteInput struct {
	Command        string `json:"command" jsonschema_description:"The shell command to execute"`
	WorkingDir     string `json:"working_dir,omitempty" jsonschema_description:"Optional directory in which to run the command (default: current working directory)"`
	TimeoutSeconds int    `json:"timeout_seconds,omitempty" jsonschema_description:"Command timeout in seconds (default: 60)"`
}

// ExecuteOutput result of shell command execution.
type ExecuteOutput struct {
	Command  string `json:"command"`
	Output   string `json:"output"`
	ExitCode int    `json:"exit_code"`
	Success  bool   `json:"success"`
	TimedOut bool   `json:"timed_out"`
}

// NewExecutionTools constructs the execution tool group.
func NewExecutionTools(backend filesystem.Shell) []tool.BaseTool {
	executeTool := Must(utils.InferTool(
		"execute",
		"Execute a shell command with timeout and working directory support, capturing stdout and stderr",
		func(ctx context.Context, input *ExecuteInput) (*ExecuteOutput, error) {
			if input.Command == "" {
				return nil, fmt.Errorf("command is required for execute tool")
			}

			timeout := 60 * time.Second
			if input.TimeoutSeconds > 0 {
				timeout = time.Duration(input.TimeoutSeconds) * time.Second
			}
			execCtx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			cmd := exec.CommandContext(execCtx, "/bin/sh", "-c", input.Command)
			if input.WorkingDir != "" {
				cmd.Dir = input.WorkingDir
			}

			out, err := cmd.CombinedOutput()
			exitCode := 0
			timedOut := false

			if err != nil {
				if execCtx.Err() == context.DeadlineExceeded {
					timedOut = true
					exitCode = -1
				} else if exitErr, ok := err.(*exec.ExitError); ok {
					exitCode = exitErr.ExitCode()
				} else {
					exitCode = 1
				}
			}

			outStr := strings.TrimSpace(string(out))
			return &ExecuteOutput{
				Command:  input.Command,
				Output:   outStr,
				ExitCode: exitCode,
				Success:  exitCode == 0 && !timedOut,
				TimedOut: timedOut,
			}, nil
		},
	))

	return []tool.BaseTool{executeTool}
}
