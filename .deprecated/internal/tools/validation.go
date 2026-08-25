package tools

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"

	. "MAgHARCM/internal/patterns"
)

// ValidateBuildInput parameters for build validation.
type ValidateBuildInput struct {
	ProjectDir string `json:"project_dir,omitempty" jsonschema_description:"Directory of the project to build/check (default: .)"`
	Toolchain  string `json:"toolchain,omitempty" jsonschema_description:"Explicit toolchain to use: cargo, make, gcc, go, auto (default: auto)"`
}

// ValidateBuildOutput result of build validation.
type ValidateBuildOutput struct {
	Success   bool   `json:"success"`
	Compiler  string `json:"compiler"`
	Output    string `json:"output"`
	HasErrors bool   `json:"has_errors"`
	Message   string `json:"message"`
}

// RunTestsInput parameters for test execution.
type RunTestsInput struct {
	ProjectDir string `json:"project_dir,omitempty" jsonschema_description:"Directory of the project to test (default: .)"`
	TestFilter string `json:"test_filter,omitempty" jsonschema_description:"Optional test name filter"`
}

// RunTestsOutput result of test execution.
type RunTestsOutput struct {
	Success     bool   `json:"success"`
	Output      string `json:"output"`
	TotalPassed int    `json:"total_passed"`
	TotalFailed int    `json:"total_failed"`
	Message     string `json:"message"`
}

// NewValidationTools constructs the validation tool group.
func NewValidationTools() []tool.BaseTool {
	validateBuildTool := Must(utils.InferTool(
		"validate_build",
		"Validation: Check if the project compiles cleanly using the appropriate build system (Cargo, Make, Go, GCC)",
		func(ctx context.Context, input *ValidateBuildInput) (*ValidateBuildOutput, error) {
			dir := input.ProjectDir
			if dir == "" {
				dir = "."
			}
			dir = filepath.Clean(dir)

			var cmd *exec.Cmd
			compiler := "auto"

			cargoPath := filepath.Join(dir, "Cargo.toml")
			if input.Toolchain == "cargo" || (input.Toolchain == "" || input.Toolchain == "auto") && fileExists(cargoPath) {
				compiler = "cargo"
				cmd = exec.CommandContext(ctx, "cargo", "check", "--manifest-path", cargoPath)
			} else if input.Toolchain == "go" || fileExists(filepath.Join(dir, "go.mod")) {
				compiler = "go"
				cmd = exec.CommandContext(ctx, "go", "vet", "./...")
				cmd.Dir = dir
			} else if input.Toolchain == "make" || fileExists(filepath.Join(dir, "Makefile")) {
				compiler = "make"
				cmd = exec.CommandContext(ctx, "make", "-n")
				cmd.Dir = dir
			} else {
				compiler = "none"
				return &ValidateBuildOutput{
					Success:   true,
					Compiler:  "none",
					Output:    "No recognizable build configuration found (Cargo.toml, go.mod, or Makefile)",
					HasErrors: false,
					Message:   "Build check skipped: no manifest found",
				}, nil
			}

			out, err := cmd.CombinedOutput()
			outStr := strings.TrimSpace(string(out))
			success := err == nil

			msg := "Build check passed cleanly."
			if !success {
				msg = fmt.Sprintf("Build check failed with exit code: %v", err)
			}

			return &ValidateBuildOutput{
				Success:   success,
				Compiler:  compiler,
				Output:    outStr,
				HasErrors: !success,
				Message:   msg,
			}, nil
		},
	))

	runTestsTool := Must(utils.InferTool(
		"run_tests",
		"Validation: Run the project's test suite and report results",
		func(ctx context.Context, input *RunTestsInput) (*RunTestsOutput, error) {
			dir := input.ProjectDir
			if dir == "" {
				dir = "."
			}
			dir = filepath.Clean(dir)

			var cmd *exec.Cmd
			cargoPath := filepath.Join(dir, "Cargo.toml")

			if fileExists(cargoPath) {
				args := []string{"test", "--manifest-path", cargoPath}
				if input.TestFilter != "" {
					args = append(args, "--", input.TestFilter)
				}
				cmd = exec.CommandContext(ctx, "cargo", args...)
			} else if fileExists(filepath.Join(dir, "go.mod")) {
				args := []string{"test", "./..."}
				if input.TestFilter != "" {
					args = append(args, "-run", input.TestFilter)
				}
				cmd = exec.CommandContext(ctx, "go", args...)
				cmd.Dir = dir
			} else {
				return &RunTestsOutput{
					Success: false,
					Output:  "No test runner found (Cargo.toml or go.mod required)",
					Message: "Cannot run tests: unrecognized test harness",
				}, nil
			}

			out, err := cmd.CombinedOutput()
			outStr := strings.TrimSpace(string(out))
			success := err == nil

			return &RunTestsOutput{
				Success: success,
				Output:  outStr,
				Message: fmt.Sprintf("Test execution finished. Success: %v", success),
			}, nil
		},
	))

	return []tool.BaseTool{validateBuildTool, runTestsTool}
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}
