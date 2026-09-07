// Package pipestate hosts the cross-cutting wire types threaded through
// every multi-agent pipeline node: the State struct (passed between agent
// Run methods), the SchemaVersioned interface, the generic DocumentWrapper
// helper, and the Task struct (user-supplied translation specification).
//
// Backlink: [[1.0.0 ADR-2026-09-07-Sprint-Conventions]] rule ADR-C-014
// (Locality of Behaviour for producer artifact structs).
//
// pipestate is a leaf package: it imports only the Go standard library so
// the broader dependency graph stays acyclic. The State struct's artifact
// fields (AnalyzerOutput, PlanningOutput, TranslatedProject,
// ValidationReport, SpecMinerInvariants, ArchaeologyReport) are stored as
// `any` and populated by the producer agent files declared canonically in
// internal/agents/. Consumers in agent files that need strongly-typed
// field access do a typed assertion (e.g.
// state.AnalyzerOutput.(agents.AnalyzerOutput)) because pipestate cannot
// import the agents package to obtain the concrete type.
//
// This trade-off is the durable resolution to the import-cycle blocker
// described at the top of internal/compiletime/state.go's prior revision
// (Sprint 2026-09-23 commit 89904f6 and Sprint 2026-09-26 subagent D
// abort). Before this refactor, compiletime.State referenced the artifact
// structs AND ValidationReport.CompilationStatus() returned
// compiletime.CompilationStatus, forcing the return type to move with
// the method receiver. That created a cycle: compiletime → agents
// (ValidationReport) → compiletime (State). Moving State out of
// compiletime into pipestate — where it can stay cycle-free by storing
// artifact fields opaquely — breaks the loop while preserving Locality
// of Behaviour for the artifacts (canonical declarations live next to
// each producer's algorithm in internal/agents/).
//
// JSON round-trip note: when a State is round-tripped through
// encoding/json (e.g. checkpoint resume), the `any` fields decode to
// map[string]interface{} rather than the original concrete type. The
// agents.DecodeStateArtifacts helper rehydrates each field into its
// canonical producer type after a checkpoint load.
package pipestate

import (
	"fmt"
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

// -------------------------------------------------------------------------
// Pipeline State — the wire type threaded through every agent Run().
// Artifact fields are stored as `any` so pipestate stays a leaf (no import
// back to internal/agents/ or internal/compiletime/). Producer agent files
// declare the concrete artifact types in their canonical home (per
// ADR-C-014). Consumers within a context that knows the concrete type
// type-assert (state.X.(agents.Y)). After a JSON round-trip (checkpoint
// resume), agents.DecodeStateArtifacts rehydrates each `any` field back
// into its canonical producer type.
// -------------------------------------------------------------------------

// State is the shared context passed between every multi-agent pipeline
// node. Each artifact field is owned by exactly one producer agent
// (algorithm in the producer file; type declared in internal/agents/).
type State struct {
	Task                Task                `json:"task"`
	AnalyzerOutput      any                 `json:"analyzer_output"`
	PlanningOutput      any                 `json:"planning_output"`
	TranslatedProject   any                 `json:"translated_project"`
	ValidationReport    any                 `json:"validation_report"`
	Iteration           int                 `json:"iteration"`
	MaxIterations       int                 `json:"max_iterations"`
	IsComplete          bool                `json:"is_complete"`
	PriorTestSnapshots  map[string]string   `json:"prior_test_snapshots,omitempty"`
	SpecMinerInvariants any                 `json:"spec_miner_invariants,omitempty"`
	ArchaeologyReport   any                 `json:"archaeology_report,omitempty"`
}
