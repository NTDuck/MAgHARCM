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
	Toolchain  string `json:"toolchain,omitempty" jsonschema_description:"Build toolchain override (e.g. cargo, go, make)"`
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
	ProjectDir string `json:"project_dir" jsonschema_description:"Directory of the project"`
	Language   string `json:"language,omitempty" jsonschema_description:"Target language (e.g. rust, c, go)"`
	Toolchain  string `json:"toolchain,omitempty" jsonschema_description:"Test toolchain override (e.g. cargo, go, make)"`
	TestFilter string `json:"test_filter,omitempty" jsonschema_description:"Optional filter/name of specific test to execute"`
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

// ValidateProjectBuild runs the appropriate compiler check based on language and toolchain.
func ValidateProjectBuild(ctx context.Context, projectDir, lang, toolchain string) (*ValidateBuildOutput, error) {
	if projectDir == "" {
		projectDir = "."
	}

	cleanDir := filepath.Clean(projectDir)
	langLower := strings.ToLower(lang)
	tcLower := strings.ToLower(toolchain)
	var cmd *exec.Cmd
	compiler := "unknown"

	cargoToml := filepath.Join(cleanDir, "Cargo.toml")
	goMod := filepath.Join(cleanDir, "go.mod")
	makefile := filepath.Join(cleanDir, "Makefile")

	// Determine compiler command
	if tcLower == "cargo" || langLower == "rust" || fileExists(cargoToml) {
		compiler = "cargo"
		cmd = exec.CommandContext(ctx, "cargo", "check", "--tests")
		cmd.Dir = cleanDir
	} else if tcLower == "go" || langLower == "go" || fileExists(goMod) {
		compiler = "go"
		cmd = exec.CommandContext(ctx, "go", "build", "./...")
		cmd.Dir = cleanDir
	} else if tcLower == "make" || langLower == "c" || langLower == "c++" || langLower == "cpp" || fileExists(makefile) {
		compiler = "make"
		cmd = exec.CommandContext(ctx, "make")
		cmd.Dir = cleanDir
	} else {
		return &ValidateBuildOutput{
			Success:   false,
			Compiler:  "none",
			Output:    "No recognized build configuration found (Cargo.toml, go.mod, Makefile)",
			HasErrors: true,
			Errors:    []string{"Missing build configuration"},
		}, nil
	}

	outBytes, err := cmd.CombinedOutput()
	outputStr := string(outBytes)

	var errs, warns []string
	for _, line := range strings.Split(outputStr, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "error[") || strings.HasPrefix(trimmed, "error:") || strings.Contains(trimmed, ": error:") {
			errs = append(errs, trimmed)
		} else if strings.HasPrefix(trimmed, "warning:") || strings.Contains(trimmed, ": warning:") {
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

// RunProjectTests executes unit and integration tests based on language and toolchain.
func RunProjectTests(ctx context.Context, projectDir, lang, toolchain, filter string) (*RunTestsOutput, error) {
	if projectDir == "" {
		projectDir = "."
	}
	cleanDir := filepath.Clean(projectDir)
	langLower := strings.ToLower(lang)
	tcLower := strings.ToLower(toolchain)

	cargoToml := filepath.Join(cleanDir, "Cargo.toml")
	goMod := filepath.Join(cleanDir, "go.mod")
	makefile := filepath.Join(cleanDir, "Makefile")

	// 1. Rust (Cargo)
	if tcLower == "cargo" || langLower == "rust" || fileExists(cargoToml) {
		args := []string{"test", "--lib", "--tests"}
		if filter != "" {
			args = append(args, filter)
		}
		args = append(args, "--", "--nocapture")

		cmd := exec.CommandContext(ctx, "cargo", args...)
		cmd.Dir = cleanDir
		outBytes, err := cmd.CombinedOutput()
		outputStr := string(outBytes)

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

	// 2. Go (go test)
	if tcLower == "go" || langLower == "go" || fileExists(goMod) {
		args := []string{"test", "-v", "./..."}
		if filter != "" {
			args = append(args, "-run", filter)
		}

		cmd := exec.CommandContext(ctx, "go", args...)
		cmd.Dir = cleanDir
		outBytes, err := cmd.CombinedOutput()
		outputStr := string(outBytes)

		passed := strings.Count(outputStr, "--- PASS:")
		failed := strings.Count(outputStr, "--- FAIL:")

		var failures []string
		if failed > 0 || err != nil {
			for _, line := range strings.Split(outputStr, "\n") {
				trimmed := strings.TrimSpace(line)
				if strings.HasPrefix(trimmed, "--- FAIL:") || strings.HasPrefix(trimmed, "FAIL") {
					failures = append(failures, trimmed)
				}
			}
		}

		success := err == nil && failed == 0 && passed > 0
		return &RunTestsOutput{
			Success:     success,
			Output:      outputStr,
			TotalPassed: passed,
			TotalFailed: failed,
			Failures:    failures,
			Message:     fmt.Sprintf("%d passed, %d failed", passed, failed),
		}, nil
	}

	// 3. C/C++ (make test)
	if tcLower == "make" || langLower == "c" || langLower == "c++" || fileExists(makefile) {
		cmd := exec.CommandContext(ctx, "make", "test")
		cmd.Dir = cleanDir
		outBytes, err := cmd.CombinedOutput()
		outputStr := string(outBytes)

		success := err == nil
		passed := 0
		failed := 0
		if success {
			passed = 1
		} else {
			failed = 1
		}

		return &RunTestsOutput{
			Success:     success,
			Output:      outputStr,
			TotalPassed: passed,
			TotalFailed: failed,
			Message:     fmt.Sprintf("make test executed (success=%v)", success),
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
			return ValidateProjectBuild(tCtx, input.ProjectDir, input.Language, input.Toolchain)
		})

	testTool, _ := utils.InferTool("run_tests", "Executes unit and integration tests in the project directory",
		func(ctx context.Context, input *RunTestsInput) (*RunTestsOutput, error) {
			tCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
			defer cancel()
			return RunProjectTests(tCtx, input.ProjectDir, input.Language, input.Toolchain, input.TestFilter)
		})

	return []tool.BaseTool{buildTool, testTool}
}

func fileExists(p string) bool {
	info, err := filepath.Glob(p)
	return err == nil && len(info) > 0
}
