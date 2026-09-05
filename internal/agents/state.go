package agents

import (
	"MAgHARCM/internal/compiletime"
)

// Backlink: [[1.0.0 ADR-2026-09-07-Sprint-Conventions]] (rule ADR-C-014).
//
// Per ADR-C-014, the canonical home of the State schema and the six
// pipeline artefact structs is internal/compiletime/state.go. This file
// declares type aliases so callsites that reference agents.State,
// agents.AnalyzerOutput, … continue to compile unchanged while the
// single source of truth lives in compiletime.
//
// Locality of Behaviour is preserved because each producer agent's Run()
// method still lives in its own file (analyzer.go, planning.go, …) and
// only writes its own artefact via the alias.

// State is the shared context passed between multi-agent pipeline nodes.
type State = compiletime.State

// SchemaVersioned is implemented by every artefact that carries an
// ArtifactSchemaVersion field.
type SchemaVersioned = compiletime.SchemaVersioned
