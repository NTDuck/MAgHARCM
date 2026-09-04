package agents

// Backlink: [[Methodology]] §1 Stage 5 & 6 and [[Primitives]] §NEW-PRIM-05, §NEW-PRIM-06, §NEW-PRIM-13.

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"MAgHARCM/internal/compiletime"
	"MAgHARCM/internal/logger"
	"MAgHARCM/internal/tools"
)

// FileStatus records the build/test outcome of an individual file.
type FileStatus struct {
	Path      string `json:"path"`
	Kind      string `json:"kind"` // "source" or "test"
	Compiles  bool   `json:"compiles"`
	TestPass  bool   `json:"test_pass"`
	LineCount int    `json:"line_count"`
	Error     string `json:"error,omitempty"`
}

// OptionalCheckResult is the persisted shape of an auxiliary validator primitive outcome.
type OptionalCheckResult struct {
	Name    string `json:"name"`
	Verdict string `json:"verdict"`
	Detail  string `json:"detail,omitempty"`
}

// ValidationReport is the structured report produced by the validator agent.
type ValidationReport struct {
	ArtifactSchemaVersion        string                `json:"schema_version"`
	AllSuccess                   bool                  `json:"all_success"`
	CompilationSuccess           bool                  `json:"compilation_success"`
	TestPassRate                 float64               `json:"test_pass_rate"`
	TotalTests                   int                   `json:"total_tests"`
	PassedTests                  int                   `json:"passed_tests"`
	FailedTests                  int                   `json:"failed_tests"`
	RealTests                    int                   `json:"real_tests"`
	MinRealTests                 int                   `json:"min_real_tests"`
	CompilationErrors            []string              `json:"compilation_errors"`
	TestFailures                 []string              `json:"test_failures"`
	UncoveredFunctions           []string              `json:"uncovered_functions"`
	CoverageGapReport            string                `json:"coverage_gap_report"`
	Diagnostics                  string                `json:"diagnostics"`
	PerFile                      []FileStatus          `json:"per_file"`
	IterationStart               time.Time             `json:"iteration_start,omitempty"`
	IterationWallMs              int64                 `json:"iteration_wall_ms,omitempty"`
	RemedyIterations             int                   `json:"remedy_iterations,omitempty"`
	PlateauDetected              bool                  `json:"plateau_detected,omitempty"`
	AdversarialWeakeningDetected bool                  `json:"adversarial_weakening_detected,omitempty"`
	WeakeningReasons             []string              `json:"weakening_reasons,omitempty"`
	ASTSyntaxErrors              []string              `json:"ast_syntax_errors,omitempty"`
	OptionalCheckResults         []OptionalCheckResult `json:"optional_check_results,omitempty"`
}

func (v ValidationReport) SchemaVersion() string { return v.ArtifactSchemaVersion }

func (v ValidationReport) IsAllSuccess() bool {
	return v.AllSuccess && v.CompilationSuccess && v.FailedTests == 0 && len(v.CompilationErrors) == 0 &&
		(!v.AdversarialWeakeningDetected) &&
		(v.TotalTests == 0 || v.PassedTests > 0) &&
		v.RealTests >= v.MinRealTests
}

// CompilationStatus returns the binary per-project compilation status (PASS or FAIL).
func (v ValidationReport) CompilationStatus() compiletime.CompilationStatus {
	if v.CompilationSuccess {
		return compiletime.CompilationStatusPass
	}
	return compiletime.CompilationStatusFail
}

func (v ValidationReport) String() string {
	if v.IsAllSuccess() {
		return fmt.Sprintf("Validation SUCCESS: compilation=%s, passed=%d/%d (%.1f%%), real_tests=%d (min=%d)",
			v.CompilationStatus(), v.PassedTests, v.TotalTests, v.TestPassRate, v.RealTests, v.MinRealTests)
	}
	return fmt.Sprintf("Validation INCOMPLETE: compilation=%s, passed=%d/%d (%.1f%%), compile_errs=%d, test_fails=%d, uncovered=%d\nDiagnostics:\n%s",
		v.CompilationStatus(), v.PassedTests, v.TotalTests, v.TestPassRate, len(v.CompilationErrors), len(v.TestFailures), len(v.UncoveredFunctions), v.Diagnostics)
}

// ValidatorAgent executes build and test suites, detects coverage gaps, and triggers test synthesis.
type ValidatorAgent struct {
	Model model.BaseChatModel
	// RunID identifies the current translation run; when set, a checkpoint of
	// *types.State is persisted under .artifacts/<RunID>/checkpoints/ at the
	// end of every Run (both success and error paths). Empty RunID disables
	// checkpointing — used by callers that don't want disk persistence.
	RunID string
	// Plateau applies CodaMOSA-style coverage-plateau detection
	// (NEW-PRIM-27, P49 ICSE 2023). It bounds the coverage-remedy loop:
	// when consecutive samples show marginal gain below threshold, the
	// loop terminates instead of burning the LLM budget. Initialized lazily
	// by NewValidatorAgent; tests that don't care about plateau may leave it nil.
	Plateau *PlateauDetector
	// OptionalChecks are auxiliary validation primitives executed after the
	// core cascade (AST → compile → test → weakening guard → plateau) finishes
	// successfully. Each check receives the current state, returns a label and
	// a verdict, and contributes to the validation report. Wired here so that
	// adding a new optional check (PRIM-7 verdict, PRIM-8 mock, PRIM-11 IO,
	// PRIM-12 wasm) does NOT require touching the core cascade. nil-safe.
	OptionalChecks []OptionalCheck
}

// OptionalCheck is the contract implemented by auxiliary validation primitives.
// A check is invoked only after the core cascade succeeds; its verdict is
// attached to ValidationReport.OptionalCheckResults for downstream inspection.
// Implementations MUST be idempotent — the validator may call them on every
// repair iteration until the cascade converges.
type OptionalCheck interface {
	// Name returns the primitive id (e.g. "PRIM-7-verdict-panel"). Stable across runs.
	Name() string
	// Run inspects the current state and returns a short verdict ("pass", "fail", "skipped")
	// and a human-readable detail string. Returned error indicates a transport-level
	// failure and aborts the optional-check phase (the core cascade is unaffected).
	Run(ctx context.Context, state *State) (verdict, detail string, err error)
}

// MaxRemedyIterations caps the bounded coverage-remediation loop driven by
// remedyWithPlateau. CodaMOSA itself uses maxStallLen=25 over MOSA iterations;
// we adopt a smaller ceiling because each iteration is an LLM test-synthesis
// call (much more expensive than a MOSA mutation), so 3 retries is a safe
// upper bound before declaring plateau and exiting.
const MaxRemedyIterations = 3

// NewValidatorAgent creates a ValidatorAgent instance. runID enables
// per-run checkpoint persistence; pass "" to disable checkpointing.
func NewValidatorAgent(m model.BaseChatModel, runID string) *ValidatorAgent {
	return &ValidatorAgent{Model: m, RunID: runID, Plateau: NewPlateauDetector()}
}

// Run validates the target project and produces a ValidationReport.
func (v *ValidatorAgent) Run(ctx context.Context, state *State) (*State, error) {
	defer v.checkpoint(state)
	state.Iteration++
	iterStart := time.Now()
	logger.LogAgent("Validator", "Running build and test validation on target `%s` (Iteration %d/%d)",
		state.Task.TargetDir, state.Iteration, state.MaxIterations)
	report := ValidationReport{
		ArtifactSchemaVersion: CurrentSchemaVersion,
		IterationStart:        iterStart,
	}

	// Step 1: Pre-compilation AST syntax check (NEW-PRIM-6 / GAP-08)
	syntaxErrs := v.checkASTSyntax(state.Task.TargetDir, state.Task.TargetLang)
	report.ASTSyntaxErrors = syntaxErrs
	if len(syntaxErrs) > 0 {
		logger.LogStep("AST Syntax check: %d syntax issues flagged before compilation", len(syntaxErrs))
	}

	buildRes, err := v.checkCompilation(ctx, state)
	if err != nil {
		return nil, err
	}

	report.CompilationSuccess = buildRes.Success
	report.CompilationErrors = buildRes.Errors

	if !buildRes.Success {
		report.Diagnostics = fmt.Sprintf("Compilation errors:\n%s\nCompiler Output:\n%s", strings.Join(buildRes.Errors, "\n"), buildRes.Output)
		report.PerFile = scanTargetFiles(state, buildRes.Errors)
		report.IterationWallMs = time.Since(iterStart).Milliseconds()
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
	report.RealTests = testRes.RealTests
	// MinRealTests scales with DISTINCT source-file basenames, not raw fragment
	// count. planning.go emits one fragment per AST element (e.g., 517 for
	// Sample 2 stats vs 73 source files); multiplying raw fragments by 2
	// self-inflates the threshold and makes the ceiling unreachable. Counting
	// distinct basenames matches the chunked translator's actual emit unit.
	distinctSources := len(GroupFragmentsBySourceFile(state.PlanningOutput.Fragments))
	if distinctSources == 0 {
		distinctSources = len(state.PlanningOutput.Fragments)
	}
	report.MinRealTests = max(5, distinctSources*2)
	if report.TotalTests > 0 {
		report.TestPassRate = float64(report.PassedTests) / float64(report.TotalTests) * 100.0
	}

	uncovered := v.findUncoveredFunctions(state)
	report.UncoveredFunctions = uncovered
	if (len(uncovered) > 0 || report.RealTests < report.MinRealTests) && v.Model != nil {
		v.remedyCoverageGaps(ctx, state, uncovered, &report)
	}

	report.PerFile = scanTargetFiles(state, buildRes.Errors)

	// Step 2: Adversarial weakening detection (NEW-PRIM-13 / GAP-01)
	currentTests := v.collectCurrentTestFiles(state)
	if state.PriorTestSnapshots != nil {
		weakened, reasons := v.DetectTestWeakening(state.PriorTestSnapshots, currentTests)
		report.AdversarialWeakeningDetected = weakened
		report.WeakeningReasons = reasons
		if weakened {
			logger.LogWarning("Adversarial test weakening detected: %v", reasons)
		}
	}
	state.PriorTestSnapshots = currentTests

	v.finalizeReport(&report, state, testRes.Output)
	if len(v.OptionalChecks) > 0 && report.IsAllSuccess() {
		report.OptionalCheckResults = v.runOptionalChecks(ctx, state)
	}
	report.IterationWallMs = time.Since(iterStart).Milliseconds()
	state.ValidationReport = report
	logger.LogStep("ITER[%d] comp=%v tests=%d/%d (%.1f%%) wall=%dms per-file=%d",
		state.Iteration, report.CompilationSuccess, report.PassedTests, report.TotalTests, report.TestPassRate,
		report.IterationWallMs, len(report.PerFile))
	logger.LogStep("ITER[%d] real_tests=%d min=%d (vacuous=%v)", state.Iteration, report.RealTests, report.MinRealTests, report.RealTests == 0)
	return state, nil
}

// checkpoint persists a snapshot of state if RunID is set. Errors are logged
// but never propagated — checkpointing is best-effort and must not abort the
// pipeline when the disk is full or the runID is empty.
func (v *ValidatorAgent) checkpoint(state *State) {
	if v.RunID == "" || state == nil {
		return
	}
	if path, err := Save(v.RunID, state); err != nil {
		logger.LogWarning("Validator checkpoint save failed: %v", err)
	} else {
		logger.LogStep("Validator checkpoint saved: %s", path)
	}
}

// checkCompilation verifies compilation using the target toolchain.
func (v *ValidatorAgent) checkCompilation(ctx context.Context, state *State) (*tools.ValidateBuildOutput, error) {
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
func (v *ValidatorAgent) runTestSuite(ctx context.Context, state *State) (*tools.RunTestsOutput, error) {
	testRes, err := tools.RunProjectTests(ctx, state.Task.TargetDir, state.Task.TargetLang, state.Task.Toolchain, "")
	if err != nil {
		return nil, fmt.Errorf("validator failed to execute tests: %w", err)
	}
	logger.LogTool("run_tests", "Test suite finished: %d passed, %d failed", testRes.TotalPassed, testRes.TotalFailed)
	return testRes, nil
}

// findUncoveredFunctions scans test files and matches discovered AST fragments against test bodies.
func (v *ValidatorAgent) findUncoveredFunctions(state *State) []string {
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

// remedyCoverageGaps is the public entry point for coverage-gap remediation.
// It now delegates to remedyWithPlateau, which wraps the per-iteration
// (generate tests → re-run suite) cycle in a bounded plateau-detected loop
// (CodaMOSA / NEW-PRIM-27). The wrapped inner behavior is unchanged.
func (v *ValidatorAgent) remedyCoverageGaps(ctx context.Context, state *State, uncovered []string, report *ValidationReport) {
	v.remedyWithPlateau(ctx, state, uncovered, report)
}

// remedyWithPlateau executes the coverage-remediation loop bounded by
// MaxRemedyIterations and terminated early when the PlateauDetector signals
// insufficient marginal gain across consecutive samples. Each iteration:
//
//  1. Re-discovers uncovered functions (the set may shrink as tests are added).
//  2. Calls generateAdditionalTests to prompt the LLM for supplemental tests.
//  3. Re-runs the test suite and records a CoverageSample in the plateau detector.
//  4. Breaks when IsPlateau() returns true.
//
// On exit, report.RemedyIterations and report.PlateauDetected reflect the
// observed trajectory. If v.Plateau is nil (legacy callers / unit tests that
// construct ValidatorAgent directly), a fresh detector is allocated so the
// loop is always safe.
func (v *ValidatorAgent) remedyWithPlateau(ctx context.Context, state *State, uncovered []string, report *ValidationReport) {
	detector := v.Plateau
	if detector == nil {
		detector = NewPlateauDetector()
		v.Plateau = detector
	}

	// Initial sample: cover whatever the caller observed before this loop ran.
	initialUncovered := len(uncovered)
	if len(report.UncoveredFunctions) > 0 {
		initialUncovered = len(report.UncoveredFunctions)
	}
	detector.Record(CoverageSample{
		Iteration:      0,
		TotalTests:     report.TotalTests,
		PassedTests:    report.PassedTests,
		UncoveredCount: initialUncovered,
		CapturedAt:     time.Now(),
	})

	for iter := 1; iter <= MaxRemedyIterations; iter++ {
		// Re-discover: prior iterations may have added test coverage.
		currentUncovered := v.findUncoveredFunctions(state)
		report.UncoveredFunctions = currentUncovered

		logger.LogStep("Plateau remedy iteration %d/%d: %d uncovered functions",
			iter, MaxRemedyIterations, len(currentUncovered))

		v.generateAdditionalTests(ctx, state, currentUncovered, report)

		logger.LogStep("Re-running test suite after coverage test generation")
		if reTest, err := tools.RunProjectTests(ctx, state.Task.TargetDir, state.Task.TargetLang, state.Task.Toolchain, ""); err == nil {
			report.PassedTests = reTest.TotalPassed
			report.FailedTests = reTest.TotalFailed
			report.TotalTests = reTest.TotalPassed + reTest.TotalFailed
			report.TestFailures = reTest.Failures
			report.RealTests = reTest.RealTests
			if report.TotalTests > 0 {
				report.TestPassRate = float64(report.PassedTests) / float64(report.TotalTests) * 100.0
			}
		}

		report.RemedyIterations = iter
		plateau := detector.Record(CoverageSample{
			Iteration:      iter,
			TotalTests:     report.TotalTests,
			PassedTests:    report.PassedTests,
			UncoveredCount: len(report.UncoveredFunctions),
			CapturedAt:     time.Now(),
		})
		if plateau {
			report.PlateauDetected = true
			logger.LogStep("Coverage plateau detected at iteration %d: %s", iter, detector.Summary())
			break
		}

		// Early exit: nothing left to chase.
		if len(report.UncoveredFunctions) == 0 && report.RealTests >= report.MinRealTests {
			break
		}
	}

	logger.LogStep("%s", detector.Summary())
}

// finalizeReport evaluates convergence criteria and sets milestone diagnostics.
func (v *ValidatorAgent) finalizeReport(report *ValidationReport, state *State, testOutput string) {
	report.AllSuccess = report.CompilationSuccess && len(report.CompilationErrors) == 0 && report.FailedTests == 0 && report.RealTests >= report.MinRealTests
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

func (v *ValidatorAgent) generateAdditionalTests(ctx context.Context, state *State, uncovered []string, report *ValidationReport) {
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

	prompt, err := renderPromptTemplate("validator_coverage", validatorCoveragePromptTemplate, map[string]any{
		"TargetLang":         state.Task.TargetLang,
		"TargetLangLower":    strings.ToLower(state.Task.TargetLang),
		"UncoveredFunctions": strings.Join(uncovered, "\n"),
		"SourceFiles":        strings.Join(sourceFilesData, "\n"),
		"TestFileRelPath":    testFileRelPath,
		"MinRealTests":       strconv.Itoa(report.MinRealTests),
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

// fileErrorPattern matches compiler error lines that reference a file path.
var fileErrorPattern = regexp.MustCompile(`(?:^|["'\s])([A-Za-z0-9_./\-]+\.[A-Za-z0-9]+):(?:\d+:\d+)?\s*(?:error|warning|note)`)

// scanTargetFiles enumerates target source + test files and joins each file to
// any compilation error mentioning it. Used for per-file observability.
func scanTargetFiles(state *State, compileErrors []string) []FileStatus {
	targetDir := state.Task.TargetDir
	statuses := []FileStatus{}
	seen := map[string]bool{}

	// First: in-memory translated project files (most authoritative).
	for relPath, content := range state.TranslatedProject.Files {
		kind := "source"
		if strings.HasPrefix(relPath, "tests/") || strings.HasSuffix(relPath, "_test."+strings.ToLower(state.Task.TargetLang)) {
			kind = "test"
		}
		statuses = append(statuses, FileStatus{
			Path:      relPath,
			Kind:      kind,
			Compiles:  true, // assume yes; failure is signaled by error grep below
			LineCount: strings.Count(content, "\n") + 1,
		})
		seen[relPath] = true
	}

	// Then: any files on disk that aren't in memory (e.g. skeleton files written but never edited).
	_ = filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(targetDir, path)
		if rel == "" || seen[rel] || strings.HasPrefix(rel, "target/") || strings.HasPrefix(rel, "target\\") || strings.HasPrefix(rel, ".") {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".rs" && ext != ".toml" && ext != ".go" && ext != ".java" && ext != ".py" && ext != ".c" && ext != ".h" && ext != ".cc" && ext != ".cpp" && ext != ".js" && ext != ".ts" {
			return nil
		}
		data, _ := os.ReadFile(path)
		kind := "source"
		if strings.HasPrefix(rel, "tests/") || strings.HasSuffix(rel, "_test."+strings.ToLower(state.Task.TargetLang)) {
			kind = "test"
		}
		statuses = append(statuses, FileStatus{
			Path:      rel,
			Kind:      kind,
			Compiles:  true,
			LineCount: strings.Count(string(data), "\n") + 1,
		})
		return nil
	})

	// Now attribute compile errors to files (best-effort grep).
	for i := range statuses {
		for _, errLine := range compileErrors {
			if fileErrorPattern.MatchString(errLine) && strings.Contains(errLine, statuses[i].Path) {
				statuses[i].Compiles = false
				statuses[i].Error = errLine
				break
			}
		}
	}
	return statuses
}

// checkASTSyntax scans target source files for tree-sitter syntax errors before running compiler (NEW-PRIM-6 / GAP-08).
func (v *ValidatorAgent) checkASTSyntax(targetDir, targetLang string) []string {
	var syntaxErrors []string
	_ = filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext == ".rs" || ext == ".go" || ext == ".java" || ext == ".py" || ext == ".c" || ext == ".cpp" {
			structOut, parseErr := tools.ParseFileStructure(path)
			if parseErr != nil {
				syntaxErrors = append(syntaxErrors, fmt.Sprintf("%s: parse error: %v", filepath.Base(path), parseErr))
			} else if structOut != nil && len(structOut.Elements) == 0 && info.Size() > 200 {
				syntaxErrors = append(syntaxErrors, fmt.Sprintf("%s: 0 AST elements extracted from %d bytes", filepath.Base(path), info.Size()))
			}
		}
		return nil
	})
	return syntaxErrors
}

// countAssertions counts assertion calls in test source code (NEW-PRIM-13 / GAP-01).
func countAssertions(code string) int {
	re := regexp.MustCompile(`(?i)\b(?:assert|expect|should|require)\w*\s*[\(!]`)
	return len(re.FindAllString(code, -1))
}

// DetectTestWeakening detects whether tests were modified, deleted, or weakened during repair (NEW-PRIM-13 / GAP-01).
func (v *ValidatorAgent) DetectTestWeakening(prevTests, currentTests map[string]string) (bool, []string) {
	if len(prevTests) == 0 {
		return false, nil
	}
	var reasons []string
	for file, prevContent := range prevTests {
		currContent, exists := currentTests[file]
		if !exists {
			reasons = append(reasons, fmt.Sprintf("Test file `%s` was removed", file))
			continue
		}
		prevAsserts := countAssertions(prevContent)
		currAsserts := countAssertions(currContent)
		if currAsserts < prevAsserts {
			reasons = append(reasons, fmt.Sprintf("Test file `%s`: assertions reduced from %d to %d", file, prevAsserts, currAsserts))
		}
	}
	return len(reasons) > 0, reasons
}

// collectCurrentTestFiles extracts a snapshot of all test files in the target project.
func (v *ValidatorAgent) collectCurrentTestFiles(state *State) map[string]string {
	tests := make(map[string]string)
	for path, content := range state.TranslatedProject.Files {
		if strings.Contains(path, "test") || strings.HasSuffix(path, "_test.go") || strings.HasSuffix(path, "_test.rs") {
			tests[path] = content
		}
	}
	return tests
}

// runOptionalChecks executes each registered OptionalCheck against the
// current state. Errors are recorded as "fail" with the error message as
// detail; the core cascade result is unchanged. Used to wire PRIM-7 verdict
// panel, PRIM-8 mock validator, PRIM-11 implementation-agnostic I/O tests,
// and PRIM-12 wasm oracle into the repair loop without touching the core
// cascade. Order follows v.OptionalChecks (callers register in the order
// they want results recorded).
func (v *ValidatorAgent) runOptionalChecks(ctx context.Context, state *State) []OptionalCheckResult {
	results := make([]OptionalCheckResult, 0, len(v.OptionalChecks))
	for _, check := range v.OptionalChecks {
		if check == nil {
			continue
		}
		verdict, detail, err := check.Run(ctx, state)
		if err != nil {
			verdict = "fail"
			detail = "check error: " + err.Error()
			logger.LogWarning("OptionalCheck %s failed: %v", check.Name(), err)
		}
		results = append(results, OptionalCheckResult{
			Name:    check.Name(),
			Verdict: verdict,
			Detail:  detail,
		})
		logger.LogStep("OptionalCheck %s -> %s (%s)", check.Name(), verdict, truncateForLog(detail, 120))
	}
	return results
}

// truncateForLog keeps optional-check detail lines short for the log stream.
func truncateForLog(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
