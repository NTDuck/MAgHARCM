package compiletime

// Backlink: [[1.0.0 ADR-2026-09-07-Sprint-Conventions]] (rule ADR-C-014).
//
// This file is the canonical home for the schema-versioned inter-agent
// pipeline context. The State struct and its six artifact structs live here
// as concrete types. Per-agent producer files in internal/agents/ declare
// `type Foo = compiletime.Foo` aliases so callsites continue to compile
// unchanged while the single source of truth moves to compiletime.
//
// Per ADR-C-014:
//
//   - State + SchemaVersioned: HERE (canonical).
//   - Per-artifact structs: HERE (canonical) — they are pipeline artefacts,
//     not agent-internal state. Locality of Behaviour is preserved because
//     each agent's Run() method still lives in its own file and writes its
//     own artefact through the alias.
//   - Agents reference via `type Foo = compiletime.Foo` alias in their
//     respective files.

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

// -------------------------------------------------------------------------
// Analyzer artefacts (PRIM-1 / NEW-PRIM-20 / NEW-PRIM-21 / NEW-PRIM-26)
// -------------------------------------------------------------------------

// DocumentWrapper keeps both structured data and markdown representation.
type DocumentWrapper[T any] struct {
	ArtifactSchemaVersion string `json:"schema_version"`
	Data                  T      `json:"data"`
	RawMarkdown           string `json:"raw_markdown"`
}

// SchemaVersion implements SchemaVersioned.
func (d DocumentWrapper[T]) SchemaVersion() string { return d.ArtifactSchemaVersion }

// SourceProjectResearch represents the research document produced by AnalyzerAgent.
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

// ThirdPartyLibraryAnalysis represents the library analysis document.
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

// TargetProjectDesign represents the design document produced by AnalyzerAgent.
type TargetProjectDesign struct {
	Overview                string   `json:"overview"`
	TranslationRequirements string   `json:"translation_requirements"`
	SourceFilesToTranslate  []string `json:"source_files_to_translate"`
	ModuleStructure         string   `json:"module_structure"`
	ErrorHandling           string   `json:"error_handling"`
	ThirdPartyLibraries     []string `json:"third_party_libraries"`
	RawDocument             string   `json:"raw_document"`
}

// AnalyzerOutput aggregates research, library mapping, and architectural
// design documents. Producer: AnalyzerAgent (internal/agents/analyzer.go).
type AnalyzerOutput struct {
	ArtifactSchemaVersion string                                     `json:"schema_version"`
	Research              DocumentWrapper[SourceProjectResearch]     `json:"research"`
	Library               DocumentWrapper[ThirdPartyLibraryAnalysis] `json:"library"`
	Design                DocumentWrapper[TargetProjectDesign]       `json:"design"`
}

// SchemaVersion implements SchemaVersioned.
func (a AnalyzerOutput) SchemaVersion() string { return a.ArtifactSchemaVersion }

// -------------------------------------------------------------------------
// Planning artefacts (PRIM-3 / NEW-PRIM-01 / NEW-PRIM-02)
// -------------------------------------------------------------------------

// PlanStep represents a single step in Part A or Part B of the
// implementation plan.
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

// PlanningOutput captures AST fragments, symbol mappings, generated
// skeletons, and translation steps. Producer: PlanningAgent.
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
// Translator artefact (PRIM-6 / NEW-PRIM-23 / NEW-PRIM-24 / NEW-PRIM-25)
// -------------------------------------------------------------------------

// TranslatedProject contains the files written or edited in the target
// repository. Producer: TranslatorAgent.
type TranslatedProject struct {
	ArtifactSchemaVersion string            `json:"schema_version"`
	Files                 map[string]string `json:"files"` // relative_path -> code_content
}

// SchemaVersion implements SchemaVersioned.
func (t TranslatedProject) SchemaVersion() string { return t.ArtifactSchemaVersion }

// -------------------------------------------------------------------------
// Validator artefact (PRIM-5 / NEW-PRIM-13)
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

// ValidationReport is the structured report produced by the validator
// agent. Producer: ValidatorAgent.
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
// (PASS or FAIL). Per ADR-C-009 there is no partial compilation rate.
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
// SpecMiner artefact (PRIM-4)
// -------------------------------------------------------------------------

// SpecMinerInvariants represents the dynamic invariants recovered by
// SpecMiner. Producer: SpecMiner.
type SpecMinerInvariants struct {
	AllocSizes         []int              `json:"alloc_sizes"`
	PointerNullability map[string]bool    `json:"pointer_nullability"`
	AliasingPairs      [][2]string        `json:"aliasing_pairs"`
	LifetimeRanges     map[string][2]int  `json:"lifetime_ranges"`
	BranchCoverage     map[string]float64 `json:"branch_coverage"`
}

// -------------------------------------------------------------------------
// Archaeologist artefact (PRIM-14)
// -------------------------------------------------------------------------

// NamingFinding records a legacy-identifier forensic finding produced by
// the Archaeologist.
type NamingFinding struct {
	Style, File, Token, Suggestion string
	Line                           int
}

// ArchaeologyReport is the structured output of the archaeologist agent.
// Producer: Archaeologist.
type ArchaeologyReport struct {
	BoundaryMap         []string            // module / file / function boundaries
	TimeCapsuleCommands []string            // shell commands that reproduce the original build / test env
	ChurnHotspots       []string            // top-N churned files from git history
	NamingForensics     []NamingFinding     // legacy-identifier findings (Hungarian, m_, ALL_CAPS, …)
	ConceptMap          map[string][]string // co-occurrence concept clusters
}

// -------------------------------------------------------------------------
// Pipeline State (the wire type threaded through every agent Run())
// -------------------------------------------------------------------------

// State is the shared context passed between multi-agent pipeline nodes.
// Each artifact field is owned by one producer agent (Locality of Behaviour).
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
