package types

import (
	"fmt"
	"time"
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

// ValidationReport is the structured report produced by the validator agent.
// A run is AllSuccess ONLY when compilation passes AND every real test passes
// AND the MinRealTests gate is satisfied.
type ValidationReport struct {
	AllSuccess         bool         `json:"all_success"`
	CompilationSuccess bool         `json:"compilation_success"`
	TestPassRate       float64      `json:"test_pass_rate"`
	TotalTests         int          `json:"total_tests"`
	PassedTests        int          `json:"passed_tests"`
	FailedTests        int          `json:"failed_tests"`
	RealTests          int          `json:"real_tests"`
	MinRealTests       int          `json:"min_real_tests"`
	CompilationErrors  []string     `json:"compilation_errors"`
	TestFailures       []string     `json:"test_failures"`
	UncoveredFunctions []string     `json:"uncovered_functions"`
	CoverageGapReport  string       `json:"coverage_gap_report"`
	Diagnostics        string       `json:"diagnostics"`
	PerFile            []FileStatus `json:"per_file"`
	IterationStart     time.Time    `json:"iteration_start,omitempty"`
	IterationWallMs    int64        `json:"iteration_wall_ms,omitempty"`
	// RemedyIterations counts coverage-remedy iterations executed by the
	// plateau-bounded loop (NEW-PRIM-27 / CodaMOSA). 0 means the loop never
	// entered; MaxRemedyIterations means the loop exhausted without plateau.
	RemedyIterations int `json:"remedy_iterations,omitempty"`
	// PlateauDetected is true when the PlateauDetector signaled insufficient
	PlateauDetected bool `json:"plateau_detected,omitempty"`
	// AdversarialWeakeningDetected is true when test weakening (e.g. dropped assertions,
	// emptied test functions, relaxed predicates) is detected across repair iterations (NEW-PRIM-13).
	AdversarialWeakeningDetected bool     `json:"adversarial_weakening_detected,omitempty"`
	WeakeningReasons             []string `json:"weakening_reasons,omitempty"`
	// ASTSyntaxErrors contains parse/syntax errors found before compilation (NEW-PRIM-6 / GAP-08).
	ASTSyntaxErrors []string `json:"ast_syntax_errors,omitempty"`
}

// HasUncoveredFunctions determines if any discovered AST functions lack test assertions.
func (v *ValidationReport) HasUncoveredFunctions() bool {
	return len(v.UncoveredFunctions) > 0
}

// IsAllSuccess returns true ONLY when the validator finalized report.AllSuccess
// with the MinRealTests gate satisfied. The struct field is the single source
// of truth — duplicate field-level checks here caused the gate to be ignored
// by graph.go / main.go (which call this method).
func (v *ValidationReport) IsAllSuccess() bool {
	return v.AllSuccess
}

// String provides a human-readable summary of the validation report.
func (v *ValidationReport) String() string {
	if v.IsAllSuccess() {
		return fmt.Sprintf("Validation SUCCESS: %d/%d tests passed (100.0%% pass rate)", v.PassedTests, v.TotalTests)
	}
	return fmt.Sprintf("Validation INCOMPLETE: compilation=%v, passed=%d/%d (%.1f%%), compile_errs=%d, test_fails=%d, uncovered=%d\nDiagnostics:\n%s",
		v.CompilationSuccess, v.PassedTests, v.TotalTests, v.TestPassRate, len(v.CompilationErrors), len(v.TestFailures), len(v.UncoveredFunctions), v.Diagnostics)
}
