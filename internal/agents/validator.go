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

// ValidatorAgent executes build and test suites, detects coverage gaps, and triggers test synthesis.
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
	logger.LogAgent("Validator", "Validating target project `%s` (Iteration %d/%d)",
		state.Task.TargetDir, state.Iteration, state.MaxIterations)

	buildRes, err := v.checkCompilation(ctx, state)
	if err != nil {
		return nil, err
	}

	report := types.ValidationReport{
		CompilationSuccess: buildRes.Success,
		CompilationErrors:  buildRes.Errors,
	}

	if !buildRes.Success {
		report.Diagnostics = fmt.Sprintf("Compilation errors:\n%s", strings.Join(buildRes.Errors, "\n"))
		state.ValidationReport = report
		return state, nil
	}

	testRes, err := v.runTestSuite(ctx, state)
	if err != nil {
		return nil, err
	}

	report.TotalTests = testRes.TotalPassed + testRes.TotalFailed
	report.PassedTests = testRes.TotalPassed
	report.FailedTests = testRes.TotalFailed
	report.TestFailures = testRes.Failures
	if report.TotalTests > 0 {
		report.TestPassRate = float64(report.PassedTests) / float64(report.TotalTests) * 100.0
	}

	uncovered := v.findUncoveredFunctions(state)
	report.UncoveredFunctions = uncovered
	if len(uncovered) > 0 && testRes.Success && testRes.TotalPassed == 0 && v.Model != nil {
		v.remedyCoverageGaps(ctx, state, uncovered, &report)
	}

	v.finalizeReport(&report, state, testRes.Output)
	state.ValidationReport = report
	return state, nil
}

// checkCompilation verifies compilation using the target toolchain.
func (v *ValidatorAgent) checkCompilation(ctx context.Context, state *types.State) (*tools.ValidateBuildOutput, error) {
	buildRes, err := tools.ValidateProjectBuild(ctx, state.Task.TargetDir, state.Task.TargetLang, state.Task.Toolchain)
	if err != nil {
		return nil, fmt.Errorf("validator failed to run build check: %w", err)
	}

	if !buildRes.Success {
		logger.LogTool("validate_build", "Compilation FAILED with %d errors:\n%s",
			len(buildRes.Errors), strings.Join(buildRes.Errors, "\n"))
	} else {
		logger.LogTool("validate_build", "Compilation SUCCESS: clean build without errors")
	}
	return buildRes, nil
}
// runTestSuite executes the target test suite and records test metrics.
func (v *ValidatorAgent) runTestSuite(ctx context.Context, state *types.State) (*tools.RunTestsOutput, error) {
	testRes, err := tools.RunProjectTests(ctx, state.Task.TargetDir, state.Task.TargetLang, state.Task.Toolchain, "")
	if err != nil {
		return nil, fmt.Errorf("validator failed to execute tests: %w", err)
	}
	logger.LogTool("run_tests", "Test suite finished: %d passed, %d failed", testRes.TotalPassed, testRes.TotalFailed)
	return testRes, nil
}

// findUncoveredFunctions scans test files and matches discovered AST fragments against test bodies.
func (v *ValidatorAgent) findUncoveredFunctions(state *types.State) []string {
	var testFilesContent []string
	_ = filepath.Walk(filepath.Join(state.Task.TargetDir, "tests"), func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			if data, err := os.ReadFile(path); err == nil {
				testFilesContent = append(testFilesContent, string(data))
			}
		}
		return nil
	})
	allTestsCode := strings.Join(testFilesContent, "\n")

	var uncovered []string
	for _, frag := range state.PlanningOutput.Fragments {
		parts := strings.Split(frag, ":")
		if len(parts) == 2 {
			fnName := parts[1]
			if !strings.Contains(allTestsCode, fnName) {
				uncovered = append(uncovered, fnName)
			}
		}
	}
	return uncovered
}

// remedyCoverageGaps prompts the model for supplemental tests and re-executes tests.
func (v *ValidatorAgent) remedyCoverageGaps(ctx context.Context, state *types.State, uncovered []string, report *types.ValidationReport) {
	logger.LogStep("Coverage Gap Analysis: %d functions uncovered and 0 tests executed, prompting Reasoning Model for tests", len(uncovered))
	v.generateAdditionalTests(ctx, state, uncovered)

	logger.LogStep("Re-running test suite after coverage test generation")
	if reTest, err := tools.RunProjectTests(ctx, state.Task.TargetDir, state.Task.TargetLang, state.Task.Toolchain, ""); err == nil {
		report.PassedTests = reTest.TotalPassed
		report.FailedTests = reTest.TotalFailed
		report.TotalTests = reTest.TotalPassed + reTest.TotalFailed
		report.TestFailures = reTest.Failures
		if report.TotalTests > 0 {
			report.TestPassRate = float64(report.PassedTests) / float64(report.TotalTests) * 100.0
		}
	}
}

// finalizeReport evaluates convergence criteria and sets milestone diagnostics.
func (v *ValidatorAgent) finalizeReport(report *types.ValidationReport, state *types.State, testOutput string) {
	report.AllSuccess = report.CompilationSuccess && len(report.CompilationErrors) == 0 && report.FailedTests == 0 && report.PassedTests > 0
	if report.AllSuccess {
		report.Diagnostics = fmt.Sprintf("All %d tests passed successfully! Codebase compiled cleanly without errors.", report.PassedTests)
		state.IsComplete = true
		logger.LogAgent("Validator", "Validation SUCCESS: %s", report.Diagnostics)
	} else {
		report.Diagnostics = fmt.Sprintf("Tests failed or 0 tests passed: %d passed, %d failed.\nFailures:\n%s\nOutput:\n%s",
			report.PassedTests, report.FailedTests, strings.Join(report.TestFailures, "\n"), testOutput)
		logger.LogAgent("Validator", "Validation INCOMPLETE: %d passed, %d failed", report.PassedTests, report.FailedTests)
	}
}

// generateAdditionalTests requests supplemental test cases from the reasoning model.
func (v *ValidatorAgent) generateAdditionalTests(ctx context.Context, state *types.State, uncovered []string) {
	var sourceFilesData []string
	for relPath, content := range state.TranslatedProject.Files {
		if !strings.HasPrefix(relPath, "tests/") {
			sourceFilesData = append(sourceFilesData, fmt.Sprintf("=== File: %s ===\n%s\n", relPath, content))
		}
	}

	testFileRelPath := "tests/integration_tests.rs"
	for relPath := range state.TranslatedProject.Files {
		if strings.HasPrefix(relPath, "tests/") {
			testFileRelPath = relPath
			break
		}
	}

	prompt, err := renderPrompt("validator_coverage.md", map[string]any{
		"TargetLang":         state.Task.TargetLang,
		"TargetLangLower":    strings.ToLower(state.Task.TargetLang),
		"UncoveredFunctions": strings.Join(uncovered, "\n"),
		"SourceFiles":        strings.Join(sourceFilesData, "\n"),
		"TestFileRelPath":    testFileRelPath,
	})
	if err != nil {
		logger.LogError("Failed to render validator coverage prompt: %v", err)
		return
	}
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
				logger.LogTool("write_file", "Updated tests in `%s` (%d bytes)", relPath, len(cleaned))
			}
		}
	}
}
