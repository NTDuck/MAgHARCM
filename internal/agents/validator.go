package agents

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"

	"MAgHARCM/internal/logger"
	"MAgHARCM/internal/tools"
	"MAgHARCM/internal/types"
)

// ValidatorAgent performs compilation checks, test execution, coverage-gap analysis, and test generation (§3.5).
type ValidatorAgent struct {
	Model model.BaseChatModel
}

// NewValidatorAgent creates a ValidatorAgent instance.
func NewValidatorAgent(m model.BaseChatModel) *ValidatorAgent {
	return &ValidatorAgent{Model: m}
}

// Run validates the target project and produces a ValidationReport.
func (v *ValidatorAgent) Run(ctx context.Context, state *types.State) (*types.State, error) {
	state.Iteration++
	logger.LogAgent("Validator", "Validating target project %s (Iteration %d/%d)",
		state.Task.TargetDir, state.Iteration, state.MaxIterations)
	state.Log("[Validator] Validating target project %s (Iteration %d/%d)", state.Task.TargetDir, state.Iteration, state.MaxIterations)

	// 1. Validate build compilation (§3.5.1)
	logger.LogStep("Running compiler check: cargo check --tests...")
	buildRes, err := tools.ValidateProjectBuild(ctx, state.Task.TargetDir, state.Task.TargetLang)
	if err != nil {
		return nil, fmt.Errorf("validator failed to run build check: %w", err)
	}

	report := types.ValidationReport{
		CompilationSuccess: buildRes.Success,
		CompilationErrors:  buildRes.Errors,
	}

	if !buildRes.Success {
		report.AllSuccess = false
		report.Diagnostics = fmt.Sprintf("Compilation failed with errors:\n%s\nOutput:\n%s",
			strings.Join(buildRes.Errors, "\n"), buildRes.Output)
		state.ValidationReport = report
		logger.LogTool("validate_build", "Compilation FAILED with %d errors:\n%s",
			len(buildRes.Errors), strings.Join(buildRes.Errors, "\n"))
		state.Log("[Validator] Build check FAILED with %d compilation errors", len(buildRes.Errors))
		return state, nil
	}
	logger.LogTool("validate_build", "Compilation SUCCESS: clean build without errors")

	// 2. Run unit tests (§3.5.1)
	logger.LogStep("Executing test suite: cargo test -- --nocapture...")
	testRes, err := tools.RunProjectTests(ctx, state.Task.TargetDir, "")
	if err != nil {
		return nil, fmt.Errorf("validator failed to execute tests: %w", err)
	}

	report.TotalTests = testRes.TotalPassed + testRes.TotalFailed
	report.PassedTests = testRes.TotalPassed
	report.FailedTests = testRes.TotalFailed
	report.TestFailures = testRes.Failures
	if report.TotalTests > 0 {
		report.TestPassRate = float64(report.PassedTests) / float64(report.TotalTests) * 100.0
	}
	logger.LogTool("run_tests", "Test suite finished: %d passed, %d failed", report.PassedTests, report.FailedTests)

	// 3. Coverage-Guided Test Generation (§3.5.2)
	var testFilesContent []string
	_ = filepath.Walk(filepath.Join(state.Task.TargetDir, "tests"), func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if data, err := os.ReadFile(path); err == nil {
			testFilesContent = append(testFilesContent, string(data))
		}
		return nil
	})
	allTestsCode := strings.Join(testFilesContent, "\n")

	var uncovered []string
	for _, frag := range state.PlanningOutput.Fragments {
		parts := strings.Split(frag, ":")
		if len(parts) == 2 {
			funcName := parts[1]
			if !strings.Contains(allTestsCode, funcName) {
				uncovered = append(uncovered, frag)
			}
		}
	}
	report.UncoveredFunctions = uncovered

	if len(uncovered) > 0 && testRes.Success && testRes.TotalPassed == 0 && v.Model != nil {
		logger.LogStep("Coverage Gap Analysis: %d functions uncovered and 0 tests executed, prompting Reasoning Model for tests...", len(uncovered))
		v.generateAdditionalTests(ctx, state, uncovered)
		// Re-run tests after adding generated tests
		logger.LogStep("Re-running test suite after coverage test generation...")
		if reTest, err := tools.RunProjectTests(ctx, state.Task.TargetDir, ""); err == nil {
			report.TotalTests = reTest.TotalPassed + reTest.TotalFailed
			report.PassedTests = reTest.TotalPassed
			report.FailedTests = reTest.TotalFailed
			report.TestFailures = reTest.Failures
			if report.TotalTests > 0 {
				report.TestPassRate = float64(report.PassedTests) / float64(report.TotalTests) * 100.0
			}
			logger.LogTool("run_tests", "Updated test suite: %d passed, %d failed", report.PassedTests, report.FailedTests)
		}
	}
	report.AllSuccess = buildRes.Success && len(report.CompilationErrors) == 0 && report.PassedTests > 0
	if report.AllSuccess {
		report.Diagnostics = fmt.Sprintf("Compilation succeeded and %d/%d tests passed (%.1f%% pass rate > 0%%)",
			report.PassedTests, report.TotalTests, report.TestPassRate)
		state.IsComplete = true
		logger.LogAgent("Validator", "Validation SUCCESS: %s", report.Diagnostics)
		state.Log("[Validator] Validation SUCCESS: %s", report.Diagnostics)
	} else {
		report.Diagnostics = fmt.Sprintf("Tests failed or 0 tests passed: %d passed, %d failed.\nFailures:\n%s\nOutput:\n%s",
			report.PassedTests, report.FailedTests, strings.Join(report.TestFailures, "\n"), testRes.Output)
		logger.LogAgent("Validator", "Validation INCOMPLETE: %d passed, %d failed", report.PassedTests, report.FailedTests)
		state.Log("[Validator] Validation INCOMPLETE: %d passed, %d failed", report.PassedTests, report.FailedTests)
	}
	state.ValidationReport = report
	return state, nil
}

func (v *ValidatorAgent) generateAdditionalTests(ctx context.Context, state *types.State, uncovered []string) {
	var targetFilesContent []string
	_ = filepath.Walk(filepath.Join(state.Task.TargetDir, "src"), func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if data, err := os.ReadFile(path); err == nil {
			targetFilesContent = append(targetFilesContent, fmt.Sprintf("=== File: %s ===\n%s\n", filepath.Base(path), string(data)))
		}
		return nil
	})

	// Find existing test file name
	testRelPath := "tests/integration_tests.rs"
	_ = filepath.Walk(filepath.Join(state.Task.TargetDir, "tests"), func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".rs") {
			rel, _ := filepath.Rel(state.Task.TargetDir, path)
			testRelPath = rel
		}
		return nil
	})

	prompt := fmt.Sprintf(`You are the Validator Agent in ReCodeAgent performing Coverage-Guided Test Generation (§3.5.2).
The following functions or modules are uncovered or need more test assertions:
%s

Target Source Files:
%s

Generate additional comprehensive unit test cases in the target test framework to thoroughly exercise all uncovered functions and edge cases.
Output the complete updated test file:

FILE: %s
`+"```rust"+`
// Complete updated test suite with additional test assertions
`+"```"+`
`, strings.Join(uncovered, "\n"), strings.Join(targetFilesContent, "\n"), testRelPath)

	resp, err := v.Model.Generate(ctx, []*schema.Message{
		schema.SystemMessage("You are an expert test engineer writing thorough unit tests."),
		schema.UserMessage(prompt),
	})
	if err == nil && resp != nil {
		files := parseAllFileMarkers(resp.Content)
		for relPath, testCode := range files {
			if strings.HasPrefix(relPath, "tests/") {
				cleaned := tools.CleanCodeContent(testCode)
				testPath := filepath.Join(state.Task.TargetDir, relPath)
				_ = os.WriteFile(testPath, []byte(cleaned), 0644)
				state.TranslatedProject.Files[relPath] = cleaned
				logger.LogTool("write_file", "Updated tests in %s (%d bytes)", relPath, len(cleaned))
			}
		}
	}
}
