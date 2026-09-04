package agents

import (
	"MAgHARCM/internal/compiletime"
)


// CurrentSchemaVersion is the schema version stamped on every role artifact.
const CurrentSchemaVersion = compiletime.CurrentSchemaVersion

// SchemaVersioned is implemented by every artifact that carries a SchemaVersion field.
type SchemaVersioned interface {
	SchemaVersion() string
}
// State is the shared context passed between multi-agent pipeline nodes.
// Each field is a structured artifact produced and owned by one agent,
// preserving Locality of Behaviour and high cohesion.
type State struct {
	Task                compiletime.Task    `json:"task"`
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
