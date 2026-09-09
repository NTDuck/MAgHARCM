package agents

// Typed structured-output contracts for each LLM-calling agent. Every agent
// that consumes model output has a struct here that mirrors its artifact in
// internal/compiletime/state.go. The structs are the JSON schema the model is
// asked to fill in via a single forced tool call — replacing the brittle
// text/regex parsing path the agents used to perform on raw markdown.
//
// Design choices:
//   - Schemas live in the agents package (where the algorithms run), not in
//     compiletime, so the producer agent owns both the artifact and the
//     schema that produces it. compiletime keeps the rich artifact type
//     with full method receivers (IsAllSuccess, etc.) and accepts the
//     typed fields directly from these contracts.
//   - JSON tags are the schema source of truth; do not introduce a parallel
//     struct definition.
//   - Required fields are encoded as non-pointer types so the inferred
//     JSON-Schema marks them required; optional metadata uses omitempty
//     or pointer types.

import (
	"MAgHARCM/internal/compiletime"
)

// AnalyzerSchema is the typed output of the analyzer agent. The model fills
// every field directly; the post-processing path used to slice this out of
// a raw markdown document.
type AnalyzerSchema struct {
	Overview           string                       `json:"overview"`
	DirectoryStructure string                       `json:"directory_structure"`
	StructsInterfaces  string                       `json:"structs_and_interfaces"`
	DataModels         string                       `json:"data_models"`
	ErrorHandling      string                       `json:"error_handling"`
	Dependencies       []string                     `json:"dependencies"`
	Libraries          []AnalyzerLibraryEntrySchema `json:"libraries"`
	Design             AnalyzerDesignSchema         `json:"design"`
}

// AnalyzerLibraryEntrySchema mirrors compiletime.LibraryMapping; it lives here
// because it is part of the analyzer's tool-call contract, not a top-level
// pipeline artifact.
type AnalyzerLibraryEntrySchema struct {
	SourceLibrary   string `json:"source_library"`
	TargetLibrary   string `json:"target_library"`
	Overview        string `json:"overview"`
	Usage           string `json:"usage"`
	Recommendations string `json:"recommendations"`
}

// AnalyzerDesignSchema mirrors the fields the planner consumes from
// compiletime.TargetProjectDesign. Migration strategy and rationale come
// from SelectAndTryStrategies and are merged in by the analyzer post-process.
type AnalyzerDesignSchema struct {
	Overview                string   `json:"overview"`
	TranslationRequirements string   `json:"translation_requirements"`
	SourceFilesToTranslate  []string `json:"source_files_to_translate"`
	ModuleStructure         string   `json:"module_structure"`
	ErrorHandling           string   `json:"error_handling"`
	ThirdPartyLibraries     []string `json:"third_party_libraries"`
}

// toAnalyzerOutput projects the typed schema into the canonical compiletime
// artifact with strategy fields populated from the strategy selector.
func (a AnalyzerSchema) toAnalyzerOutput(strategy compiletime.StrategyKind, rationale, rawMarkdown string) compiletime.AnalyzerOutput {
	libs := make([]compiletime.LibraryMapping, 0, len(a.Libraries))
	for _, l := range a.Libraries {
		libs = append(libs, compiletime.LibraryMapping{
			SourceLibrary:   l.SourceLibrary,
			TargetLibrary:   l.TargetLibrary,
			Overview:        l.Overview,
			Usage:           l.Usage,
			Recommendations: l.Recommendations,
		})
	}
	return compiletime.AnalyzerOutput{
		ArtifactSchemaVersion: compiletime.CurrentSchemaVersion,
		Research: compiletime.DocumentWrapper[compiletime.SourceProjectResearch]{
			ArtifactSchemaVersion: compiletime.CurrentSchemaVersion,
			Data: compiletime.SourceProjectResearch{
				Overview:           a.Overview,
				DirectoryStructure: a.DirectoryStructure,
				StructsInterfaces:  a.StructsInterfaces,
				DataModels:         a.DataModels,
				ErrorHandling:      a.ErrorHandling,
				Dependencies:       a.Dependencies,
				MigrationStrategy:  string(strategy),
				StrategyRationale:  rationale,
				RawDocument:        rawMarkdown,
			},
			RawMarkdown: rawMarkdown,
		},
		Library: compiletime.DocumentWrapper[compiletime.ThirdPartyLibraryAnalysis]{
			ArtifactSchemaVersion: compiletime.CurrentSchemaVersion,
			Data: compiletime.ThirdPartyLibraryAnalysis{
				Libraries:   libs,
				RawDocument: rawMarkdown,
			},
			RawMarkdown: rawMarkdown,
		},
		Design: compiletime.DocumentWrapper[compiletime.TargetProjectDesign]{
			ArtifactSchemaVersion: compiletime.CurrentSchemaVersion,
			Data: compiletime.TargetProjectDesign{
				Overview:                a.Design.Overview,
				TranslationRequirements: a.Design.TranslationRequirements,
				SourceFilesToTranslate:  a.Design.SourceFilesToTranslate,
				ModuleStructure:         a.Design.ModuleStructure,
				ErrorHandling:           a.Design.ErrorHandling,
				ThirdPartyLibraries:     a.Design.ThirdPartyLibraries,
				RawDocument:             rawMarkdown,
			},
			RawMarkdown: rawMarkdown,
		},
	}
}

// PlanningSchema is the typed output of the planning agent. The model emits
// the symbol name map, the skeleton file contents keyed by relative path, and
// the implementation plan (overview + ordered Part A / Part B steps).
type PlanningSchema struct {
	NameMapping   map[string]string             `json:"name_mapping"`
	SkeletonFiles map[string]string             `json:"skeleton_files"` // relative_path -> full file content
	Overview      string                        `json:"overview"`
	PartA         []compiletime.PlanStep        `json:"part_a"`
	PartB         []compiletime.PlanStep        `json:"part_b"`
}

// TranslatorSchema is the typed output of every translator path
// (single-shot translate, repair, and chunked translateFragment). The model
// emits a map of relative path -> full file content; the agent persists them
// via syncFilesToDisk.
type TranslatorSchema struct {
	Files map[string]string `json:"files"` // relative_path -> full file content
}

// ValidatorCoverageSchema is the typed output of the validator's coverage-
// guided test generation pass. Only the test files are consumed; production
// files emitted by the model are dropped.
type ValidatorCoverageSchema struct {
	Files map[string]string `json:"files"`
}

// JudgeSchema is the typed output of one equivalence-judge vote (PRIM-7).
// The model emits a single canonical verdict token plus a free-text rationale.
type JudgeSchema struct {
	Verdict   string `json:"verdict"`   // "EQUIVALENT" or "NOT_EQUIVALENT"
	Rationale string `json:"rationale"`
}

// RoleFlipSchema is the typed output of the communicative de-hallucination
// reviewer (PRIM-25). The model emits either a defect claim or an explicit
// acceptance; both are captured as structured fields.
type RoleFlipSchema struct {
	Defect   bool   `json:"defect"`
	Reason   string `json:"reason"`
	RetryHint string `json:"retry_hint,omitempty"`
}
