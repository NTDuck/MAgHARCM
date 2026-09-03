// Package types holds the shared contract types that flow through the
// translation pipeline. State is the single bag the Eino graph nodes pass
// between one another; the per-domain files in this package hold the
// artifacts each agent produces.
//
// Cross-package moves are deliberately avoided: every artifact type in
// State is referenced by at least two packages (producer + consumer) and
// moving them would create import cycles. The per-file split inside this
// package gives the Locality of Behaviour benefit without that risk.
package types

// State is the shared context passed between Eino graph nodes.
type State struct {
	Task               TranslationTask   `json:"task"`
	AnalyzerOutput     AnalyzerOutput    `json:"analyzer_output"`
	PlanningOutput     PlanningOutput    `json:"planning_output"`
	TranslatedProject  TranslatedProject `json:"translated_project"`
	ValidationReport   ValidationReport  `json:"validation_report"`
	Iteration          int               `json:"iteration"`
	MaxIterations      int               `json:"max_iterations"`
	IsComplete         bool              `json:"is_complete"`
	PriorTestSnapshots map[string]string `json:"prior_test_snapshots,omitempty"`
}
