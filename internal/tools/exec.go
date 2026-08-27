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

	"MAgHARCM/internal/languages"
)

// ValidateBuildInput parameters for build validation.
type ValidateBuildInput struct {
	ProjectDir string `json:"project_dir" jsonschema_description:"Directory of the project to check/build"`
	Language   string `json:"language,omitempty" jsonschema_description:"Target language (e.g. rust, c, go, python, java)"`
	Toolchain  string `json:"toolchain,omitempty" jsonschema_description:"Build toolchain override (e.g. cargo, go, make, cmake, maven, gradle, pytest)"`
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
	Language   string `json:"language,omitempty" jsonschema_description:"Target language (e.g. rust, c, go, python, java)"`
	Toolchain  string `json:"toolchain,omitempty" jsonschema_description:"Test toolchain override (e.g. cargo, go, make, cmake, maven, gradle, pytest)"`
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

// ValidateProjectBuild runs the appropriate compiler check based on language registry and toolchain.
func ValidateProjectBuild(ctx context.Context, projectDir, lang, toolchain string) (*ValidateBuildOutput, error) {
	if projectDir == "" {
		projectDir = "."
	}

	cleanDir := filepath.Clean(projectDir)
	reg := languages.GetRegistry()
	var cmd *exec.Cmd
	compiler := "unknown"

	// 1. Language-specified toolchain lookup
	if lang != "" {
		if spec, ok := reg.FindByName(lang); ok {
			for _, tc := range spec.Toolchains {
				if toolchain != "" && !strings.EqualFold(tc.Name, toolchain) {
					continue
				}
				matched := false
				for _, bf := range tc.BuildFiles {
					if fileExists(filepath.Join(cleanDir, bf)) {
						matched = true
						break
					}
				}
				if !matched {
					for _, ext := range spec.Extensions {
						if fileExists(filepath.Join(cleanDir, "*"+ext)) {
							matched = true
							break
						}
					}
				}
				if matched && len(tc.BuildCommand) > 0 {
					compiler = tc.Name
					cmd = exec.CommandContext(ctx, tc.BuildCommand[0], tc.BuildCommand[1:]...)
					cmd.Dir = cleanDir
					break
				}
			}
		}
	}

	// 2. Dynamic directory inspection across all registered language toolchains
	if cmd == nil {
		for _, name := range []string{"rust", "go", "c", "cpp", "java", "python", "typescript", "kotlin"} {
			if spec, ok := reg.FindByName(name); ok {
				for _, tc := range spec.Toolchains {
					matched := false
					for _, bf := range tc.BuildFiles {
						if fileExists(filepath.Join(cleanDir, bf)) {
							matched = true
							break
						}
					}
					if !matched {
						for _, ext := range spec.Extensions {
							if fileExists(filepath.Join(cleanDir, "*"+ext)) {
								matched = true
								break
							}
						}
					}
					if matched && len(tc.BuildCommand) > 0 {
						compiler = tc.Name
						cmd = exec.CommandContext(ctx, tc.BuildCommand[0], tc.BuildCommand[1:]...)
						cmd.Dir = cleanDir
						break
					}
				}
			}
			if cmd != nil {
				break
			}
		}
	}

	if cmd == nil {
		return &ValidateBuildOutput{
			Success:   false,
			Compiler:  "none",
			Output:    "No recognized build configuration found for project",
			HasErrors: true,
			Errors:    []string{"Missing build configuration"},
		}, nil
	}

	outBytes, err := cmd.CombinedOutput()
	outputStr := string(outBytes)

	var errs, warns []string
	for _, line := range strings.Split(outputStr, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "error[") || strings.HasPrefix(trimmed, "error:") || strings.Contains(trimmed, ": error:") || strings.HasPrefix(trimmed, "FAILED") {
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

// RunProjectTests executes unit and integration tests based on language registry and toolchain.
func RunProjectTests(ctx context.Context, projectDir, lang, toolchain, filter string) (*RunTestsOutput, error) {
	if projectDir == "" {
		projectDir = "."
	}
	cleanDir := filepath.Clean(projectDir)
	reg := languages.GetRegistry()
	var cmd *exec.Cmd
	tcName := "unknown"

	// 1. Language-specified test toolchain lookup
	if lang != "" {
		if spec, ok := reg.FindByName(lang); ok {
			for _, tc := range spec.Toolchains {
				if toolchain != "" && !strings.EqualFold(tc.Name, toolchain) {
					continue
				}
				matched := false
				for _, bf := range tc.BuildFiles {
					if fileExists(filepath.Join(cleanDir, bf)) {
						matched = true
						break
					}
				}
				if !matched {
					for _, ext := range spec.Extensions {
						if fileExists(filepath.Join(cleanDir, "*"+ext)) {
							matched = true
							break
						}
					}
				}
				if matched && len(tc.TestCommand) > 0 {
					tcName = tc.Name
					args := append([]string{}, tc.TestCommand[1:]...)
					if filter != "" {
						args = append(args, filter)
					}
					cmd = exec.CommandContext(ctx, tc.TestCommand[0], args...)
					cmd.Dir = cleanDir
					break
				}
			}
		}
	}

	// 2. Dynamic directory inspection across all registered language test runners
	if cmd == nil {
		for _, name := range []string{"rust", "go", "c", "cpp", "java", "python", "typescript", "kotlin"} {
			if spec, ok := reg.FindByName(name); ok {
				for _, tc := range spec.Toolchains {
					matched := false
					for _, bf := range tc.BuildFiles {
						if fileExists(filepath.Join(cleanDir, bf)) {
							matched = true
							break
						}
					}
					if !matched {
						for _, ext := range spec.Extensions {
							if fileExists(filepath.Join(cleanDir, "*"+ext)) {
								matched = true
								break
							}
						}
					}
					if matched && len(tc.TestCommand) > 0 {
						tcName = tc.Name
						args := append([]string{}, tc.TestCommand[1:]...)
						if filter != "" {
							args = append(args, filter)
						}
						cmd = exec.CommandContext(ctx, tc.TestCommand[0], args...)
						cmd.Dir = cleanDir
						break
					}
				}
			}
			if cmd != nil {
				break
			}
		}
	}

	if cmd == nil {
		return &RunTestsOutput{
			Success: false,
			Output:  "No test runner found for directory: " + projectDir,
			Message: "Missing test configuration",
		}, nil
	}

	outBytes, err := cmd.CombinedOutput()
	outputStr := string(outBytes)

	// Universal test result parsing
	passed := 0
	failed := 0

	passRe := regexp.MustCompile(`(?i)(?:(\d+)\s+passed|---\s*PASS:|PASSED)`)
	failRe := regexp.MustCompile(`(?i)(?:(\d+)\s+failed|---\s*FAIL:|FAILED)`)

	for _, match := range passRe.FindAllStringSubmatch(outputStr, -1) {
		if len(match) > 1 && match[1] != "" {
			if n, err := strconv.Atoi(match[1]); err == nil {
				passed += n
				continue
			}
		}
		passed++
	}
	for _, match := range failRe.FindAllStringSubmatch(outputStr, -1) {
		if len(match) > 1 && match[1] != "" {
			if n, err := strconv.Atoi(match[1]); err == nil {
				failed += n
				continue
			}
		}
		failed++
	}

	if passed == 0 && failed == 0 {
		if err == nil {
			passed = 1
		} else {
			failed = 1
		}
	}

	var failures []string
	if failed > 0 || err != nil {
		for _, line := range strings.Split(outputStr, "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.Contains(trimmed, "FAILED") || strings.Contains(trimmed, "FAIL") || strings.Contains(trimmed, "panicked at") || strings.Contains(trimmed, "error[") || strings.Contains(trimmed, "error:") {
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
		Message:     fmt.Sprintf("%s test run: %d passed, %d failed", tcName, passed, failed),
	}, nil
}

// NewExecutionTools creates language-agnostic build and test execution tools.
func NewExecutionTools() []tool.BaseTool {
	buildTool, _ := utils.InferTool("validate_build", "Validates project build and compilation using the configured or auto-detected language toolchain",
		func(ctx context.Context, input *ValidateBuildInput) (*ValidateBuildOutput, error) {
			tCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
			defer cancel()
			return ValidateProjectBuild(tCtx, input.ProjectDir, input.Language, input.Toolchain)
		})

	testTool, _ := utils.InferTool("run_tests", "Executes unit and integration test suites using the project language test runner",
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
