package compiletime

// Backlink: [[1.0.0 ADR-2026-09-07-Sprint-Conventions]] (rule ADR-C-014).
//
// This file is the canonical home for the cross-cutting pipeline types
// shared across every agent and the orchestration graph:
//   - State (the wire type threaded through every agent Run())
//   - Task (the user-supplied translation specification)
//   - All per-agent artifact structs (declaration here, algorithms in
//     each producer agent's own file under internal/agents/)
//   - SchemaVersioned interface
//   - DocumentWrapper[T] generic helper
//
// ADR-C-014 declares Locality of Behaviour for the algorithms: each
// producer agent's Run() method lives in its own file. But the SHARED
// TYPES — those referenced by graph.go, runner.go, checkpoint.go, and
// cross-package consumers — must have a single source of truth.
//
// The previous architecture tried to put State in package agents and
// re-export via type alias from compiletime, but Go's import-graph
// analysis rejected it as a cycle (compiletime → agents → compiletime).
// The previous-architecture also tried introducing internal/pipeline as
// an intermediary, which created the same cycle. The cycle-free pattern
// is: declarare State and its artifact field types HERE, in compiletime,
// which is the leaf package (it imports nothing inside MAgHARCM).
// Producer agent files import compiletime and consume the types
// directly. Method receivers live on the canonical types.

import (
	"fmt"
	"time"
)

// SchemaVersioned is implemented by every artifact that carries an
// ArtifactSchemaVersion field. The interface lets the checkpoint and
// translator pipelines detect schema drift across iterations.
type SchemaVersioned interface {
	SchemaVersion() string
}

// DocumentWrapper keeps both structured data and markdown representation.
// Generic so many different artifact types can wrap their structured
// data alongside the raw markdown the reasoning model produced.
type DocumentWrapper[T any] struct {
	ArtifactSchemaVersion string `json:"schema_version"`
	Data                  T      `json:"data"`
	RawMarkdown           string `json:"raw_markdown"`
}

// SchemaVersion implements SchemaVersioned.
func (d DocumentWrapper[T]) SchemaVersion() string { return d.ArtifactSchemaVersion }

// -------------------------------------------------------------------------
// Task (user-supplied translation specification; producer: cmd entry-points
// and config loader). Lives here so State can reference it without
// dragging in a non-leaf dependency.
// -------------------------------------------------------------------------

// Task defines the specification for a translation task. Every required
// field is populated by the configuration loader before pipeline execution.
type Task struct {
	SourceDir   string
	TargetDir   string
	SourceLang  string
	TargetLang  string
	Toolchain   string
	LSPProvider string
	RequestFile string
}

// Validate verifies that every required field of Task is populated.
func (t Task) Validate() error {
	if t.SourceDir == "" {
		return fmt.Errorf("task: source_dir is required")
	}
	if t.TargetDir == "" {
		return fmt.Errorf("task: target_dir is required")
	}
	if t.SourceLang == "" {
		return fmt.Errorf("task: source_lang is required")
	}
	if t.TargetLang == "" {
		return fmt.Errorf("task: target_lang is required")
	}
	if t.Toolchain == "" {
		return fmt.Errorf("task: toolchain is required")
	}
	if t.LSPProvider == "" {
		return fmt.Errorf("task: lsp_provider is required")
	}
	return nil
}

// MustTask returns the task or panics if invalid. Use at startup where a
// malformed task is a fatal configuration error.
func MustTask(t Task) Task {
	Must(struct{}{}, t.Validate())
	return t
}

// -------------------------------------------------------------------------
// Analyzer artefacts. Algorithm: internal/agents/analyzer.go (Run, etc.).
// -------------------------------------------------------------------------

// SourceProjectResearch is the structured research document produced by
// the Analyzer Agent.
type SourceProjectResearch struct {
	Overview           string   `json:"overview"`
	DirectoryStructure string   `json:"directory_structure"`
	StructsInterfaces  string   `json:"structs_and_interfaces"`
	DataModels         string   `json:"data_models"`
	ErrorHandling      string   `json:"error_handling"`
	Dependencies       []string `json:"dependencies"`
	MigrationStrategy  string   `json:"migration_strategy,omitempty"`
	StrategyRationale  string   `json:"strategy_rationale,omitempty"`
	RawDocument        string   `json:"raw_document"`
}

// ThirdPartyLibraryAnalysis is the library analysis document.
type ThirdPartyLibraryAnalysis struct {
	Libraries   []LibraryMapping `json:"libraries"`
	RawDocument string           `json:"raw_document"`
}

// LibraryMapping maps a source library to its target counterpart.
type LibraryMapping struct {
	SourceLibrary   string `json:"source_library"`
	TargetLibrary   string `json:"target_library"`
	Overview        string `json:"overview"`
	Usage           string `json:"usage"`
	Recommendations string `json:"recommendations"`
}

// TargetProjectDesign is the design document produced by the Analyzer.
type TargetProjectDesign struct {
	Overview                string   `json:"overview"`
	TranslationRequirements string   `json:"translation_requirements"`
	SourceFilesToTranslate  []string `json:"source_files_to_translate"`
	ModuleStructure         string   `json:"module_structure"`
	ErrorHandling           string   `json:"error_handling"`
	ThirdPartyLibraries     []string `json:"third_party_libraries"`
	RawDocument             string   `json:"raw_document"`
}

// AnalyzerOutput aggregates research, library mapping, and design documents.
type AnalyzerOutput struct {
	ArtifactSchemaVersion string                                     `json:"schema_version"`
	Research              DocumentWrapper[SourceProjectResearch]     `json:"research"`
	Library               DocumentWrapper[ThirdPartyLibraryAnalysis] `json:"library"`
	Design                DocumentWrapper[TargetProjectDesign]       `json:"design"`
}

// SchemaVersion implements SchemaVersioned.
func (a AnalyzerOutput) SchemaVersion() string { return a.ArtifactSchemaVersion }

// -------------------------------------------------------------------------
// Planning artefacts. Algorithm: internal/agents/planning.go (Run, etc.).
// -------------------------------------------------------------------------

// PlanStep is a single step in Part A or Part B of the implementation plan.
type PlanStep struct {
	ID              string `json:"id"`
	Description     string `json:"description"`
	SourceFile      string `json:"source_file,omitempty"`
	TargetFile      string `json:"target_file,omitempty"`
	Type            string `json:"type,omitempty"` // "source" or "test"
	Completed       bool   `json:"completed,omitempty"`
	StepName        string `json:"step_name,omitempty"`
	Details         string `json:"details,omitempty"`
	ReverseTopoRank int    `json:"reverse_topo_rank,omitempty"`
}

// ImplementationPlan organises code translation and test verification
// steps into ordered phases.
type ImplementationPlan struct {
	ArtifactSchemaVersion string     `json:"schema_version"`
	Overview              string     `json:"overview"`
	PartA                 []PlanStep `json:"part_a"` // Source code translation
	PartB                 []PlanStep `json:"part_b"` // Test code translation & validation
	RawPlan               string     `json:"raw_plan"`
}

// PlanningOutput captures fragments, name mappings, generated skeletons,
// and translation steps.
type PlanningOutput struct {
	ArtifactSchemaVersion string             `json:"schema_version"`
	Fragments             []string           `json:"fragments"`      // file_name:fragment_name
	NameMapping           map[string]string  `json:"name_mapping"`   // source_name -> target_name
	SkeletonFiles         map[string]string  `json:"skeleton_files"` // relative_path -> skeleton_content
	Plan                  ImplementationPlan `json:"plan"`
}

// SchemaVersion implements SchemaVersioned.
func (p PlanningOutput) SchemaVersion() string { return p.ArtifactSchemaVersion }

// -------------------------------------------------------------------------
// Translator artefact. Algorithm: internal/agents/translator.go (Run, etc.).
// -------------------------------------------------------------------------

// TranslatedProject contains the files written or edited in the target
// repository.
type TranslatedProject struct {
	ArtifactSchemaVersion string            `json:"schema_version"`
	Files                 map[string]string `json:"files"` // relative_path -> code_content
}

// SchemaVersion implements SchemaVersioned.
func (t TranslatedProject) SchemaVersion() string { return t.ArtifactSchemaVersion }

// -------------------------------------------------------------------------
// Validator artefacts. Algorithm: internal/agents/validator.go (Run, etc.).
// -------------------------------------------------------------------------

// FileStatus records the build/test outcome of an individual file.
type FileStatus struct {
	Path      string `json:"path"`
	Kind      string `json:"kind"` // "source" or "test"
	Compiles  bool   `json:"compiles"`
	TestPass  bool   `json:"test_pass"`
	LineCount int    `json:"line_count"`
	Error     string `json:"error,omitempty"`
}

// OptionalCheckResult is the persisted shape of an auxiliary validator
// primitive outcome.
type OptionalCheckResult struct {
	Name    string `json:"name"`
	Verdict string `json:"verdict"`
	Detail  string `json:"detail,omitempty"`
}

// ValidationReport is the structured report produced by the validator.
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

// SchemaVersion implements SchemaVersioned.
func (v ValidationReport) SchemaVersion() string { return v.ArtifactSchemaVersion }

// IsAllSuccess evaluates convergence criteria: compilation passes, all
// tests pass, no adversarial weakening, no empty-test-suite escape.
func (v ValidationReport) IsAllSuccess() bool {
	return v.AllSuccess && v.CompilationSuccess && v.FailedTests == 0 && len(v.CompilationErrors) == 0 &&
		(!v.AdversarialWeakeningDetected) &&
		(v.TotalTests == 0 || v.PassedTests > 0) &&
		v.RealTests >= v.MinRealTests
}

// CompilationStatus returns the binary per-project compilation status
// (Pass or Fail). Per ADR-C-009 there is no partial compilation rate.
func (v ValidationReport) CompilationStatus() CompilationStatus {
	if v.CompilationSuccess {
		return CompilationStatusPass
	}
	return CompilationStatusFail
}

// String renders a one-line summary of the validation report for log output.
func (v ValidationReport) String() string {
	if v.IsAllSuccess() {
		return fmt.Sprintf("VALIDATION OK (status=%s, tests=%d/%d pass=%.0f%%, errors=%d, failures=%d, uncovered=%d)",
			v.CompilationStatus(), v.PassedTests, v.TotalTests, v.TestPassRate*100,
			len(v.CompilationErrors), len(v.TestFailures), len(v.UncoveredFunctions))
	}
	return fmt.Sprintf("VALIDATION FAIL (status=%s, tests=%d/%d pass=%.0f%%, errors=%d, failures=%d, uncovered=%d, diagnostics=%s)",
		v.CompilationStatus(), v.PassedTests, v.TotalTests, v.TestPassRate*100,
		len(v.CompilationErrors), len(v.TestFailures), len(v.UncoveredFunctions), v.Diagnostics)
}

// -------------------------------------------------------------------------
// SpecMiner artefact. Algorithm: internal/agents/specminer.go (Recover).
// -------------------------------------------------------------------------

// SpecMinerInvariants are the dynamic invariants recovered by SpecMiner.
type SpecMinerInvariants struct {
	AllocSizes         []int              `json:"alloc_sizes"`
	PointerNullability map[string]bool    `json:"pointer_nullability"`
	AliasingPairs      [][2]string        `json:"aliasing_pairs"`
	LifetimeRanges     map[string][2]int  `json:"lifetime_ranges"`
	BranchCoverage     map[string]float64 `json:"branch_coverage"`
}

// -------------------------------------------------------------------------
// Archaeologist artefacts. Algorithm: internal/agents/archaeology.go (Investigate).
// -------------------------------------------------------------------------

// NamingFinding records a legacy-identifier forensic finding.
type NamingFinding struct {
	Style, File, Token, Suggestion string
	Line                           int
}

// ArchaeologyReport is the structured output of the archaeologist agent.
type ArchaeologyReport struct {
	BoundaryMap         []string            // module / file / function boundaries
	TimeCapsuleCommands []string            // shell commands that reproduce the original build / test env
	ChurnHotspots       []string            // top-N churned files from git history
	NamingForensics     []NamingFinding     // legacy-identifier findings (Hungarian, m_, ALL_CAPS, …)
	ConceptMap          map[string][]string // co-occurrence concept clusters
}

// -------------------------------------------------------------------------
// Pipeline State — the wire type threaded through every agent Run().
// -------------------------------------------------------------------------

// State is the shared context passed between every multi-agent pipeline
// node. Each artifact field is owned by exactly one producer agent
// (algorithm in the producer file; type declared here).
type State struct {
	Task                Task                `json:"task"`
	AnalyzerOutput      AnalyzerOutput      `json:"analyzer_output"`
	PlanningOutput      PlanningOutput      `json:"planning_output"`
	TranslatedProject   TranslatedProject   `json:"translated_project"`
	ValidationReport    ValidationReport    `json:"validation_report"`
	Iteration           int                 `json:"iteration"`
	MaxIterations       int                 `json:"max_iterations"`
	IsComplete          bool                `json:"is_complete"`
	PriorTestSnapshots  map[string]string   `json:"prior_test_snapshots,omitempty"`
	SpecMinerInvariants SpecMinerInvariants `json:"spec_miner_invariants,omitempty"`
	ArchaeologyReport   ArchaeologyReport   `json:"archaeology_report,omitempty"`
}
