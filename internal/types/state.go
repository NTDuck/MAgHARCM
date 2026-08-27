package types

import "fmt"
// TranslationTask defines the input specification for the translation pipeline.
type TranslationTask struct {
	SourceDir   string `json:"source_dir"`
	TargetDir   string `json:"target_dir"`
	SourceLang  string `json:"source_lang"`
	TargetLang  string `json:"target_lang"`
	Toolchain   string `json:"toolchain,omitempty"`
	LSPProvider string `json:"lsp_provider,omitempty"`
}

// SourceProjectResearch represents the research document produced in phase 3.2.1.
type SourceProjectResearch struct {
	Overview           string   `json:"overview"`
	DirectoryStructure string   `json:"directory_structure"`
	StructsInterfaces  string   `json:"structs_and_interfaces"`
	DataModels         string   `json:"data_models"`
	ErrorHandling      string   `json:"error_handling"`
	Dependencies       []string `json:"dependencies"`
	RawDocument        string   `json:"raw_document"`
}

// ThirdPartyLibraryAnalysis represents the library analysis document produced in phase 3.2.2.
type ThirdPartyLibraryAnalysis struct {
	Libraries   []LibraryMapping `json:"libraries"`
	RawDocument string           `json:"raw_document"`
}

// LibraryMapping details how a source library maps to a target library.
type LibraryMapping struct {
	SourceLibrary   string `json:"source_library"`
	TargetLibrary   string `json:"target_library"`
	Overview        string `json:"overview"`
	Usage           string `json:"usage"`
	Recommendations string `json:"recommendations"`
}

// TargetProjectDesign represents the design document produced in phase 3.2.3.
type TargetProjectDesign struct {
	Overview                string   `json:"overview"`
	TranslationRequirements string   `json:"translation_requirements"`
	SourceFilesToTranslate  []string `json:"source_files_to_translate"`
	ModuleStructure         string   `json:"module_structure"`
	ErrorHandling           string   `json:"error_handling"`
	ThirdPartyLibraries     []string `json:"third_party_libraries"`
	RawDocument             string   `json:"raw_document"`
}

// AnalyzerOutput aggregates research, library mapping, and architectural design documents.
type AnalyzerOutput struct {
	Research DocumentWrapper[SourceProjectResearch]     `json:"research"`
	Library  DocumentWrapper[ThirdPartyLibraryAnalysis] `json:"library"`
	Design   DocumentWrapper[TargetProjectDesign]       `json:"design"`
}

// DocumentWrapper keeps both structured data and markdown representation.
type DocumentWrapper[T any] struct {
	Data        T      `json:"data"`
	RawMarkdown string `json:"raw_markdown"`
}

// PlanStep represents a single step in Part A or Part B of the implementation plan.
type PlanStep struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	SourceFile  string `json:"source_file"`
	TargetFile  string `json:"target_file"`
	Type        string `json:"type"` // "source" or "test"
	Completed   bool   `json:"completed"`
}

// ImplementationPlan organizes code translation and test verification steps into ordered phases.
type ImplementationPlan struct {
	Overview string     `json:"overview"`
	PartA    []PlanStep `json:"part_a"` // Source code translation
	PartB    []PlanStep `json:"part_b"` // Test code translation & validation
	RawPlan  string     `json:"raw_plan"`
}

// PlanningOutput captures AST fragments, symbol mappings, generated skeletons, and translation steps.
type PlanningOutput struct {
	Fragments     []string          `json:"fragments"`      // file_name:fragment_name
	NameMapping   map[string]string `json:"name_mapping"`   // source_name -> target_name
	SkeletonFiles map[string]string `json:"skeleton_files"` // relative_path -> skeleton_content
	Plan          ImplementationPlan `json:"plan"`
}

// TranslatedProject contains the files written or edited in the target repository.
type TranslatedProject struct {
	Files map[string]string `json:"files"` // relative_path -> code_content
}

// ValidationReport summarizes compilation checks, test results, and compiler diagnostics.
type ValidationReport struct {
	AllSuccess         bool     `json:"all_success"`
	CompilationSuccess bool     `json:"compilation_success"`
	TestPassRate       float64  `json:"test_pass_rate"`
	TotalTests         int      `json:"total_tests"`
	PassedTests        int      `json:"passed_tests"`
	FailedTests        int      `json:"failed_tests"`
	CompilationErrors  []string `json:"compilation_errors"`
	TestFailures       []string `json:"test_failures"`
	UncoveredFunctions []string `json:"uncovered_functions"`
	CoverageGapReport  string   `json:"coverage_gap_report"`
	Diagnostics        string   `json:"diagnostics"`
}

// HasUncoveredFunctions determines if any discovered AST functions lack test assertions.
func (v *ValidationReport) HasUncoveredFunctions() bool {
	return len(v.UncoveredFunctions) > 0
}

// IsAllSuccess returns true when achieving 100% compilation success and 100% test pass rate without errors or failures.
func (v *ValidationReport) IsAllSuccess() bool {
	return v.CompilationSuccess && v.FailedTests == 0 && v.PassedTests > 0 && len(v.CompilationErrors) == 0 && len(v.TestFailures) == 0
}

// String provides a human-readable summary of the validation report.
func (v *ValidationReport) String() string {
	if v.IsAllSuccess() {
		return fmt.Sprintf("Validation SUCCESS: %d/%d tests passed (100.0%% pass rate)", v.PassedTests, v.TotalTests)
	}
	return fmt.Sprintf("Validation INCOMPLETE: compilation=%v, passed=%d/%d (%.1f%%), compile_errs=%d, test_fails=%d, uncovered=%d\nDiagnostics:\n%s",
		v.CompilationSuccess, v.PassedTests, v.TotalTests, v.TestPassRate, len(v.CompilationErrors), len(v.TestFailures), len(v.UncoveredFunctions), v.Diagnostics)
}
// State is the shared context passed between Eino graph nodes.
type State struct {
	Task             TranslationTask    `json:"task"`
	AnalyzerOutput   AnalyzerOutput     `json:"analyzer_output"`
	PlanningOutput   PlanningOutput     `json:"planning_output"`
	TranslatedProject TranslatedProject `json:"translated_project"`
	ValidationReport ValidationReport   `json:"validation_report"`
	Iteration        int                `json:"iteration"`
	MaxIterations    int                `json:"max_iterations"`
	IsComplete       bool               `json:"is_complete"`
}
