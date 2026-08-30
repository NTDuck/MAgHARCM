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

type RunTestsOutput struct {
	Success     bool     `json:"success"`
	Output      string   `json:"output"`
	TotalPassed int      `json:"total_passed"`
	TotalFailed int      `json:"total_failed"`
	RealTests   int      `json:"real_tests"`
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

	// Universal test result parsing.
	// Rust: per-test granularity from `test foo::bar ... ok|FAILED` lines; aggregate `test result: ok.` blocks for cross-check.
	// Other languages: leave RealTests = -1 (sentinel: not parsed for this language) and keep legacy regex aggregation.
	passed := 0
	failed := 0
	realTests := -1

	isRust := strings.EqualFold(lang, "rust") || strings.EqualFold(tcName, "cargo") || strings.EqualFold(tcName, "rust")

	if isRust {
		realTests = 0
		perTestRe := regexp.MustCompile(`^test\s+\S+::\S+\s+\.\.\s+(ok|FAILED|ignored)$`)
		// Cross-check: cargo prints "test result: ok. M passed; K failed; ... ignored; ..." per binary.
		// We still derive RealTests from per-test lines (excludes harness / doc-test artifacts).
		for _, line := range strings.Split(outputStr, "\n") {
			trimmed := strings.TrimSpace(line)
			m := perTestRe.FindStringSubmatch(trimmed)
			if m == nil {
				continue
			}
			switch m[1] {
			case "ok":
				passed++
				realTests++
			case "FAILED":
				failed++
				realTests++
			}
		}

		// Sanity cross-check: any aggregate `test result:` line with passed > 0 means at least one real test ran.
		// If the per-test regex missed (e.g. very old cargo format), fall back to aggregate counts.
		aggRe := regexp.MustCompile(`(?i)test result:\s*(?:ok|FAILED)\.\s*(\d+)\s+passed;\s*(\d+)\s+failed`)
		aggMatches := aggRe.FindAllStringSubmatch(outputStr, -1)
		for _, m := range aggMatches {
			if len(m) >= 3 {
				if p, perr := strconv.Atoi(m[1]); perr == nil {
					if f, ferr := strconv.Atoi(m[2]); ferr == nil {
						if p > 0 || f > 0 {
							// Only upgrade if we have nothing from per-test parsing — guards against doc-test doubling.
							if realTests == 0 {
								realTests = p + f
							}
						}
					}
				}
			}
		}

		// Zero-fn guard: if cargo printed `0 passed; 0 failed` AND no per-test granularity AND no `#[test]` fns exist in target tests,
		// confirm RealTests = 0 (not a false-positive from doc-test counting).
		if realTests == 0 {
			grepCmd := exec.CommandContext(ctx, "grep", "-r", "--include=*.rs", "-E", "fn\\s+.*test", filepath.Join(cleanDir, "tests"))
			grepCmd.Dir = cleanDir
			if grepOut, gerr := grepCmd.CombinedOutput(); gerr == nil {
				count := 0
				for _, l := range strings.Split(string(grepOut), "\n") {
					if strings.TrimSpace(l) != "" {
						count++
					}
				}
				if count == 0 {
					realTests = 0
				}
			}
		}

		// Fall back: legacy regex for `passed`/`failed` so other harnesses / framework output still parses.
		if passed == 0 && failed == 0 && realTests > 0 {
			passed = realTests
		}
	} else {
		passRe := regexp.MustCompile(`(?i)(?:(\d+)\s+passed|---\s*PASS:|PASSED)`)
		failRe := regexp.MustCompile(`(?i)(?:(\d+)\s+failed|---\s*FAIL:|FAILED)`)

		for _, match := range passRe.FindAllStringSubmatch(outputStr, -1) {
			if len(match) > 1 && match[1] != "" {
				if n, perr := strconv.Atoi(match[1]); perr == nil {
					passed += n
					continue
				}
			}
			passed++
		}
		for _, match := range failRe.FindAllStringSubmatch(outputStr, -1) {
			if len(match) > 1 && match[1] != "" {
				if n, perr := strconv.Atoi(match[1]); perr == nil {
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
	// Vacuous pass guard for Rust: if no real functional tests were discovered AND no execution error,
	// the "passed" total must be rejected even when TotalPassed > 0 (e.g. doc-test counted).
	vacuous := false
	if isRust && realTests == 0 && err == nil {
		vacuous = true
		passed = 0
		success = false
	}

	msg := fmt.Sprintf("%s test run: %d passed, %d failed", tcName, passed, failed)
	if vacuous {
		msg = "Vacuous test pass: 0 real functional tests discovered"
	}

	return &RunTestsOutput{
		Success:     success,
		Output:      outputStr,
		TotalPassed: passed,
		TotalFailed: failed,
		RealTests:   realTests,
		Failures:    failures,
		Message:     msg,
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
