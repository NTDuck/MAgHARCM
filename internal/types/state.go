// Package types holds the State container and the input contract
// (TranslationTask) for the translation pipeline. The output artifacts
// produced by each agent live in internal/artifacts so producers and
// consumers can share them without an import cycle.
//
// Dependency direction:
//
//	internal/artifacts  (no deps)
//	         ^
//	         |
//	internal/types  --> internal/agents, internal/graph, internal/runner
package types

import "MAgHARCM/internal/artifacts"

// TranslationTask defines the input specification for the translation pipeline.
// It stays in internal/types (rather than internal/artifacts) because the YAML
// loader in internal/config populates it directly from the request file; the
// type has no "produced by agent" semantics — it is the user's request, so it
// belongs with State.
type TranslationTask struct {
	SourceDir   string `json:"source_dir"`
	TargetDir   string `json:"target_dir"`
	SourceLang  string `json:"source_lang"`
	TargetLang  string `json:"target_lang"`
	Toolchain   string `json:"toolchain,omitempty"`
	LSPProvider string `json:"lsp_provider,omitempty"`
}

// State is the shared context passed between Eino graph nodes. Each field is
// a structured artifact produced by one agent and consumed by later agents;
// the underlying types live in internal/artifacts.
type State struct {
	Task               TranslationTask             `json:"task"`
	AnalyzerOutput     artifacts.AnalyzerOutput    `json:"analyzer_output"`
	PlanningOutput     artifacts.PlanningOutput    `json:"planning_output"`
	TranslatedProject  artifacts.TranslatedProject `json:"translated_project"`
	ValidationReport   artifacts.ValidationReport  `json:"validation_report"`
	Iteration          int                         `json:"iteration"`
	MaxIterations      int                         `json:"max_iterations"`
	IsComplete         bool                        `json:"is_complete"`
	PriorTestSnapshots map[string]string           `json:"prior_test_snapshots,omitempty"`
}
