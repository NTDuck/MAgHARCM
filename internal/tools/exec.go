package tools

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

// ValidateBuildInput parameters for build validation.
type ValidateBuildInput struct {
	ProjectDir string `json:"project_dir" jsonschema_description:"Directory of the project to check/build"`
	Language   string `json:"language,omitempty" jsonschema_description:"Target language (e.g. rust, c, go)"`
}

// ValidateBuildOutput result of build validation.
type ValidateBuildOutput struct {
	Success   bool     `json:"success"`
	Compiler  string   `json:"compiler"`
	Output    string   `json:"output"`
	Errors    []string `json:"errors"`
	Warnings  []string `json:"warnings"`
	HasErrors bool     `json:"has_errors"`
}

// RunTestsInput parameters for test execution.
type RunTestsInput struct {
	ProjectDir string `json:"project_dir" jsonschema_description:"Directory of the project to test"`
	TestFilter string `json:"test_filter,omitempty" jsonschema_description:"Optional test name filter"`
}

// RunTestsOutput result of test execution.
type RunTestsOutput struct {
	Success     bool     `json:"success"`
	Output      string   `json:"output"`
	TotalPassed int      `json:"total_passed"`
	TotalFailed int      `json:"total_failed"`
	Failures    []string `json:"failures"`
	Message     string   `json:"message"`
}

// ValidateProjectBuild runs the appropriate compiler check based on project structure.
func ValidateProjectBuild(ctx context.Context, projectDir, lang string) (*ValidateBuildOutput, error) {
	if projectDir == "" {
		projectDir = "."
	}

	cleanDir := filepath.Clean(projectDir)
	var cmd *exec.Cmd
	compiler := "unknown"

	// Rust
	cargoToml := filepath.Join(cleanDir, "Cargo.toml")
	if _, err := exec.LookPath("cargo"); err == nil && fileExists(cargoToml) {
		compiler = "cargo"
		cmd = exec.CommandContext(ctx, "cargo", "check", "--tests")
		cmd.Dir = cleanDir
	} else if fileExists(filepath.Join(cleanDir, "Makefile")) {
		compiler = "make"
		cmd = exec.CommandContext(ctx, "make")
		cmd.Dir = cleanDir
	} else {
		return &ValidateBuildOutput{
			Success:   false,
			Compiler:  "none",
			Output:    "No recognized build file found (e.g. Cargo.toml, Makefile)",
			HasErrors: true,
			Errors:    []string{"Missing build configuration"},
		}, nil
	}

	outBytes, err := cmd.CombinedOutput()
	outputStr := string(outBytes)

	var errs, warns []string
	for _, line := range strings.Split(outputStr, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "error[") || strings.HasPrefix(trimmed, "error:") {
			errs = append(errs, trimmed)
		} else if strings.HasPrefix(trimmed, "warning:") {
			warns = append(warns, trimmed)
		}
	}

	success := err == nil && len(errs) == 0
	return &ValidateBuildOutput{
		Success:   success,
		Compiler:  compiler,
		Output:    outputStr,
		Errors:    errs,
		Warnings:  warns,
		HasErrors: !success,
	}, nil
}

// RunProjectTests executes unit tests in the target project.
func RunProjectTests(ctx context.Context, projectDir, filter string) (*RunTestsOutput, error) {
	if projectDir == "" {
		projectDir = "."
	}
	cleanDir := filepath.Clean(projectDir)

	cargoToml := filepath.Join(cleanDir, "Cargo.toml")
	if fileExists(cargoToml) {
		args := []string{"test", "--lib", "--tests"}
		if filter != "" {
			args = append(args, filter)
		}
		args = append(args, "--", "--nocapture")

		cmd := exec.CommandContext(ctx, "cargo", args...)
		cmd.Dir = cleanDir
		outBytes, err := cmd.CombinedOutput()
		outputStr := string(outBytes)

		// Parse Cargo test output across all test suites (unit, integration, doc-tests)
		passRe := regexp.MustCompile(`(\d+)\s+passed`)
		failRe := regexp.MustCompile(`(\d+)\s+failed`)

		passed := 0
		failed := 0

		for _, match := range passRe.FindAllStringSubmatch(outputStr, -1) {
			if len(match) > 1 {
				if n, err := strconv.Atoi(match[1]); err == nil {
					passed += n
				}
			}
		}
		for _, match := range failRe.FindAllStringSubmatch(outputStr, -1) {
			if len(match) > 1 {
				if n, err := strconv.Atoi(match[1]); err == nil {
					failed += n
				}
			}
		}

		var failures []string
		if failed > 0 || err != nil {
			for _, line := range strings.Split(outputStr, "\n") {
				trimmed := strings.TrimSpace(line)
				if strings.Contains(trimmed, "FAILED") || strings.Contains(trimmed, "panicked at") || strings.Contains(trimmed, "error[") || strings.Contains(trimmed, "error:") {
					failures = append(failures, trimmed)
				}
			}
		}

		success := err == nil && failed == 0 && (passed > 0 || strings.Contains(outputStr, "test result: ok"))
		return &RunTestsOutput{
			Success:     success,
			Output:      outputStr,
			TotalPassed: passed,
			TotalFailed: failed,
			Failures:    failures,
			Message:     fmt.Sprintf("%d passed, %d failed", passed, failed),
		}, nil
	}

	return &RunTestsOutput{
		Success: false,
		Output:  "No test runner found for directory: " + projectDir,
		Message: "Missing test configuration",
	}, nil
}

// NewExecutionTools creates execution and testing tools (§3.5).
func NewExecutionTools() []tool.BaseTool {
	buildTool, _ := utils.InferTool("validate_build", "Validates that the project compiles and checks for syntax and type errors",
		func(ctx context.Context, input *ValidateBuildInput) (*ValidateBuildOutput, error) {
			tCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
			defer cancel()
			return ValidateProjectBuild(tCtx, input.ProjectDir, input.Language)
		})

	testTool, _ := utils.InferTool("run_tests", "Executes unit and integration tests in the project directory",
		func(ctx context.Context, input *RunTestsInput) (*RunTestsOutput, error) {
			tCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
			defer cancel()
			return RunProjectTests(tCtx, input.ProjectDir, input.TestFilter)
		})

	return []tool.BaseTool{buildTool, testTool}
}

func fileExists(p string) bool {
	info, err := filepath.Glob(p)
	return err == nil && len(info) > 0
}
