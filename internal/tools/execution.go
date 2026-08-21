package tools

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk/filesystem"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"

	. "MAgHARCM/internal/patterns"
)

// ExecuteInput parameters for executing shell commands.
type ExecuteInput struct {
	Command string `json:"command" jsonschema_description:"The shell command to execute"`
}

// ExecuteOutput result of shell command execution.
type ExecuteOutput struct {
	Command  string `json:"command"`
	Output   string `json:"output"`
	ExitCode int    `json:"exit_code"`
	Success  bool   `json:"success"`
}

// NewExecutionTools constructs the execution tool group using Eino's filesystem.Shell.
func NewExecutionTools(shell filesystem.Shell) []tool.BaseTool {
	executeTool := Must(utils.InferTool(
		"execute",
		"Execute a shell command in the project environment, returning command output and exit status",
		func(ctx context.Context, input *ExecuteInput) (*ExecuteOutput, error) {
			if input.Command == "" {
				return nil, fmt.Errorf("command is required for execute tool")
			}

			resp, err := shell.Execute(ctx, &filesystem.ExecuteRequest{
				Command: input.Command,
			})
			if err != nil {
				return nil, err
			}

			exitCode := 0
			if resp.ExitCode != nil {
				exitCode = *resp.ExitCode
			}

			return &ExecuteOutput{
				Command:  input.Command,
				Output:   resp.Output,
				ExitCode: exitCode,
				Success:  exitCode == 0,
			}, nil
		},
	))

	return []tool.BaseTool{executeTool}
}
